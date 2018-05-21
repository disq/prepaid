package model

type CardStatus struct {
	ID      string  `json:"id" dynamo:"id"`
	Balance float64 `json:"balance" dynamo:"balance"`
	Blocked float64 `json:"amount" dynamo:"blocked,omitempty"`
}
