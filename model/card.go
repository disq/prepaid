package model

type CardStatus struct {
	ID               string  `json:"id" dynamo:"id"`
	AvailableBalance float64 `json:"available-balance" dynamo:"available-balance"`
	BlockedAmount    float64 `json:"blocked-amount" dynamo:"blocked-amount,omitempty"`
}
