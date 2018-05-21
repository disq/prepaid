package svc

import (
	"fmt"

	"github.com/disq/prepaid"
	"github.com/disq/prepaid/model"
	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
)

// NewCard creates a new card, loaded with initialAmount of funds.
func (se *Service) NewCard(initialAmount float64) (*model.CardStatus, error) {
	if !IsPositive(initialAmount) {
		return nil, fmt.Errorf("initialAmount should be positive")
	}

	cs := model.CardStatus{
		ID:      prepaid.UUID(),
		Balance: initialAmount,
	}

	if err := se.db.CardsTable().Put(cs).Run(); err != nil {
		return nil, errors.Wrap(err, "insertCard")
	}

	return &cs, nil
}

// CardStatus returns info about a card.
func (se *Service) CardStatus(id string) (*model.CardStatus, error) {
	var cs model.CardStatus

	if err := se.db.CardsTable().Get("id", id).One(&cs); err != nil {
		return nil, errors.Wrap(err, "getCard")
	}

	return &cs, nil
}

// CardTopup tops up a card with amt funds.
func (se *Service) CardTopup(id string, amt float64) (*model.CardStatus, error) {
	if !IsPositive(amt) {
		return nil, fmt.Errorf("amt should be positive")
	}

	var cs model.CardStatus

	if err := se.db.CardsTable().Update("id", id).If("attribute_exists(id)").Add("balance", amt).Value(&cs); err != nil {
		return nil, errors.Wrap(err, "updateCard")
	}

	return &cs, nil
}

// CardStatement returns a card statement.
func (se *Service) CardStatement(id string) (*model.CardStatement, error) {
	var cs model.CardStatement

	s, err := se.CardStatus(id)
	if err != nil {
		return nil, errors.Wrap(err, "cardStatus")
	}
	cs.Status = *s

	if err := se.db.TxTable().Get("card_id", id).Index("card_id-index").All(&cs.Transactions); err != nil && err != dynamo.ErrNotFound {
		return nil, errors.Wrap(err, "getTx")
	}

	return &cs, nil
}
