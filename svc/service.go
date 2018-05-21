package svc

import (
	"context"
	"log"

	"github.com/disq/prepaid/aws"
	"github.com/disq/prepaid/db"
	"github.com/pkg/errors"
)

// Service is our main prepaid card service.
type Service struct {
	ctx    context.Context
	logger *log.Logger

	aws *aws.AWS
	db  *db.DB
}

// New creates a new Service.
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
