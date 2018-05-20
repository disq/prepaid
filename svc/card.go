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
