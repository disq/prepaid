package svc

import (
	"github.com/disq/prepaid"
	"github.com/disq/prepaid/model"
	"github.com/pkg/errors"
)

func (se *Service) NewCard(initialAmount float64) (*model.CardStatus, error) {
	cs := model.CardStatus{
		ID:               prepaid.UUID(),
		AvailableBalance: initialAmount,
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
	var cs model.CardStatus

	if err := se.db.CardsTable().Update("id", id).If("attribute_exists(id)").Add("balance", amt).Value(&cs); err != nil {
		return nil, errors.Wrap(err, "updateCard")
	}

	return &cs, nil
}
