package svc

import (
	"context"
	"log"

	"github.com/disq/prepaid/aws"
	"github.com/disq/prepaid/db"
	"github.com/pkg/errors"
)

type Service struct {
	ctx    context.Context
	logger *log.Logger

	aws *aws.AWS
	db  *db.DB
}

func New(ctx context.Context, logger *log.Logger) (*Service, error) {
	a, err := aws.New()
	if err != nil {
		return nil, errors.Wrap(err, "aws")
	}
	d, err := db.New(a, logger)
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}

	svc := Service{
		ctx:    ctx,
		logger: logger,
		aws:    a,
		db:     d,
	}

	return &svc, nil
}
