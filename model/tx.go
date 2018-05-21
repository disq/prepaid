package model

import "time"

type TxStatus struct {
	ID       string    `json:"id" dynamo:"id"`
	CardID   string    `json:"-" dynamo:"card_id"`
	Blocked  float64   `json:"blocked" dynamo:"blocked"`
	Captured float64   `json:"captured" dynamo:"captured"`
	Refunded float64   `json:"refunded" dynamo:"refunded"`
	Expires  time.Time `json:"expires" dynamo:"expires"`
}
