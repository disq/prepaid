package svc

import (
	"fmt"

	"github.com/disq/prepaid"
	"github.com/disq/prepaid/model"
	"github.com/pkg/errors"
)

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

func (se *Service) CardStatus(id string) (*model.CardStatus, error) {
	var cs model.CardStatus

	if err := se.db.CardsTable().Get("id", id).One(&cs); err != nil {
		return nil, errors.Wrap(err, "getCard")
	}

	return &cs, nil
}

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
