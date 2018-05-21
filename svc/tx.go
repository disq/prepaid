package svc

import (
	"fmt"
	"time"

	"github.com/disq/prepaid"
	"github.com/disq/prepaid/model"
	"github.com/pkg/errors"
)

const (
	// TxExpireDuration is the time uncaptured transactions and unrefunded captures expire.
	TxExpireDuration = 5 * time.Minute
)

// NewTx creates a new transaction on card cardID, blocking amt funds for merchant.
func (se *Service) NewTx(cardID string, amt float64, merchant string) (*model.TxStatus, error) {
	if !IsPositive(amt) {
		return nil, fmt.Errorf("amt should be positive")
	}

	t := time.Now().UTC()

	ts := model.TxStatus{
		ID:       prepaid.UUID(),
		CardID:   cardID,
		Merchant: merchant,

		Blocked:   amt,
		CreatedAt: model.NullTime{Valid: true, Time: t},
		ExpiresAt: model.NullTime{Valid: true, Time: t.Add(TxExpireDuration)},
	}

	// Allocate amt in cardsTable

	if err := se.db.CardsTable().Update("id", cardID).If("attribute_exists(id) AND balance >= ?", amt).Add("blocked", amt).Add("balance", -amt).Run(); err != nil {
		return nil, errors.Wrap(err, "updateBalance")
	}

	// Put amt in Tx

	if err := se.db.TxTable().Put(ts).Run(); err != nil {
		// TODO rollback cardsTable
		se.logger.Printf("WARNING(NewTx) Orphan cardsTable: id=%v amt=%v", cardID, amt)
		return nil, errors.Wrap(err, "insertTx")
	}

	return &ts, nil
}

// TxStatus returns info about a transaction.
func (se *Service) TxStatus(id string) (*model.TxStatus, error) {
	var ts model.TxStatus

	if err := se.db.TxTable().Get("id", id).One(&ts); err != nil {
		return nil, errors.Wrap(err, "getTx")
	}

	return &ts, nil
}

// TxReverse reverses part of a transaction, unblocking amt funds.
func (se *Service) TxReverse(id string, amt float64) (*model.TxStatus, error) {
	var ts model.TxStatus

	if !IsPositive(amt) {
		return nil, fmt.Errorf("amt should be positive")
	}

	if err := se.db.TxTable().Update("id", id).If("attribute_exists(id) AND blocked >= ?", amt).Add("blocked", -amt).Value(&ts); err != nil {
		return nil, errors.Wrap(err, "updateTx")
	}

	if err := se.db.CardsTable().Update("id", ts.CardID).If("attribute_exists(id) AND blocked >= ?", amt).Add("blocked", -amt).Add("balance", amt).Run(); err != nil {
		// FIXME rollback txTable
		se.logger.Printf("WARNING(TxReverse) Orphan txTable: id=%v amt=%v", id, amt)
		return nil, errors.Wrap(err, "updateBalance")
	}

	return &ts, nil
}

// TxCapture captures part of a transaction, spending amt funds from card.
func (se *Service) TxCapture(id string, amt float64) (*model.TxStatus, error) {
	var ts model.TxStatus

	if !IsPositive(amt) {
		return nil, fmt.Errorf("amt should be positive")
	}

	t := model.NullTime{Valid: true, Time: time.Now().UTC()}

	if err := se.db.TxTable().Update("id", id).If("attribute_exists(id) AND blocked >= ?", amt).Add("blocked", -amt).Add("captured", amt).Set("captured_at", t).Value(&ts); err != nil {
		return nil, errors.Wrap(err, "updateTx")
	}

	if err := se.db.CardsTable().Update("id", ts.CardID).If("attribute_exists(id) AND blocked >= ?", amt).Add("blocked", -amt).Run(); err != nil {
		// FIXME rollback txTable
		se.logger.Printf("WARNING(TxCapture) Orphan txTable: id=%v amt=%v", id, amt)
		return nil, errors.Wrap(err, "updateBalance")
	}

	return &ts, nil
}

// TxRefund refunds part of a fully-captured transaction, refunding amt funds to card.
func (se *Service) TxRefund(id string, amt float64) (*model.TxStatus, error) {
	var ts model.TxStatus

	if !IsPositive(amt) {
		return nil, fmt.Errorf("amt should be positive")
	}

	t := time.Now().UTC()
	nt := model.NullTime{Valid: true, Time: t}

	if err := se.db.TxTable().Update("id", id).If("attribute_exists(id) AND blocked = ? AND captured >= ? AND expires > ?", 0, amt, t.Unix()).Add("captured", -amt).Add("refunded", amt).Set("refunded_at", nt).Value(&ts); err != nil {
		return nil, errors.Wrap(err, "updateTx")
	}

	if err := se.db.CardsTable().Update("id", ts.CardID).If("attribute_exists(id)", amt).Add("balance", amt).Run(); err != nil {
		// FIXME rollback txTable
		se.logger.Printf("WARNING(TxRefund) Orphan txTable: id=%v amt=%v", id, amt)
		return nil, errors.Wrap(err, "updateBalance")
	}

	return &ts, nil
}

// TxCleanup scans expired transactions and refunds blocked amount to prepaid cards.
func (se *Service) TxCleanup() error {

	var res []model.TxStatus

	if err := se.db.TxTable().Scan().All(&res); err != nil {
		return err
	}

	for _, tx := range res {
		t := time.Now().UTC()

		if tx.Blocked == 0 || !tx.ExpiresAt.Valid || tx.ExpiresAt.Time.After(t) {
			continue
		}

		_, err := se.TxReverse(tx.ID, tx.Blocked)
		if err != nil {
			se.logger.Printf("Error cleaning up %v: %v", tx.ID, err)
		}
	}

	return nil
}
