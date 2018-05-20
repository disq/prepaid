package db

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/disq/prepaid/aws"
	"github.com/guregu/dynamo"
)

type DB struct {
	*dynamo.DB

	aws    *aws.AWS
	logger *log.Logger

	cardsTable string
	txTable    string
}

const (
	envCardsTable = "CARDS_TABLE"
	envTxTable    = "TX_TABLE"
)

func New(aw *aws.AWS, logger *log.Logger) (*DB, error) {
	ct := strings.TrimSpace(os.Getenv(envCardsTable))
	if ct == "" {
		return nil, fmt.Errorf("Specify cards table in env %q", envCardsTable)
	}
	tt := strings.TrimSpace(os.Getenv(envTxTable))
	if tt == "" {
		return nil, fmt.Errorf("Specify tx table in env %q", envTxTable)
	}

	return &DB{
		DB: dynamo.New(aw.Ses),

		aws:    aw,
		logger: logger,

		cardsTable: ct,
		txTable:    tt,
	}, nil
}

func (d *DB) CardsTable() dynamo.Table {
	return d.DB.Table(d.cardsTable)
}

func (d *DB) TxTable() dynamo.Table {
	return d.DB.Table(d.txTable)
}
