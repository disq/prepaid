package svc

import (
	"fmt"
	"time"

	"github.com/disq/prepaid"
	"github.com/disq/prepaid/model"
	"github.com/pkg/errors"
)

const (
	TxExpireDuration = 12 * time.Hour
)

func (se *Service) NewTx(cardID string, amt float64, merchant string) (*model.TxStatus, error) {
	if !IsPositive(amt) {
		return nil, fmt.Errorf("amt should be positive")
	}

	ts := model.TxStatus{
		ID:      prepaid.UUID(),
		CardID:  cardID,
		Blocked: amt,
		Expires: time.Now().UTC().Add(TxExpireDuration),
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

func (se *Service) TxStatus(id string) (*model.TxStatus, error) {
	var ts model.TxStatus

	if err := se.db.TxTable().Get("id", id).One(&ts); err != nil {
		return nil, errors.Wrap(err, "getTx")
	}

	return &ts, nil
}

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

func (se *Service) TxCapture(id string, amt float64) (*model.TxStatus, error) {
	var ts model.TxStatus

	if !IsPositive(amt) {
		return nil, fmt.Errorf("amt should be positive")
	}

	if err := se.db.TxTable().Update("id", id).If("attribute_exists(id) AND blocked >= ?", amt).Add("blocked", -amt).Add("captured", amt).Value(&ts); err != nil {
		return nil, errors.Wrap(err, "updateTx")
	}

	if err := se.db.CardsTable().Update("id", ts.CardID).If("attribute_exists(id) AND blocked >= ?", amt).Add("blocked", -amt).Add("balance", -amt).Run(); err != nil {
		// FIXME rollback txTable
		se.logger.Printf("WARNING(TxCapture) Orphan txTable: id=%v amt=%v", id, amt)
		return nil, errors.Wrap(err, "updateBalance")
	}

	return &ts, nil
}

func (se *Service) TxRefund(id string, amt float64) (*model.TxStatus, error) {
	var ts model.TxStatus

	if !IsPositive(amt) {
		return nil, fmt.Errorf("amt should be positive")
	}

	if err := se.db.TxTable().Update("id", id).If("attribute_exists(id) AND blocked = ? AND captured >= ?", 0, amt).Add("captured", -amt).Add("refunded", amt).Value(&ts); err != nil {
		return nil, errors.Wrap(err, "updateTx")
	}

	if err := se.db.CardsTable().Update("id", ts.CardID).If("attribute_exists(id) AND blocked >= ?", amt).Add("blocked", -amt).Add("balance", amt).Run(); err != nil {
		// FIXME rollback txTable
		se.logger.Printf("WARNING(TxRefund) Orphan txTable: id=%v amt=%v", id, amt)
		return nil, errors.Wrap(err, "updateBalance")
	}

	return &ts, nil
}
