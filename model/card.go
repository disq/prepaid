package model

// CardStatus contains info about a prepaid card.
type CardStatus struct {
	ID      string  `json:"id" dynamo:"id"`
	Balance float64 `json:"balance" dynamo:"balance"`
	Blocked float64 `json:"blocked" dynamo:"blocked,omitempty"`
}
