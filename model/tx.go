package model

// TxStatus contains info about a transaction.Null
type TxStatus struct {
	ID         string   `json:"id" dynamo:"id"`
	CardID     string   `json:"-" dynamo:"card_id"`
	Merchant   string   `json:"merchant" dynamo:"merchant"`
	Blocked    float64  `json:"blocked" dynamo:"blocked"`
	Captured   float64  `json:"captured" dynamo:"captured"`
	Refunded   float64  `json:"refunded" dynamo:"refunded"`
	CreatedAt  NullTime `json:"created_at" dynamo:"created_at"`
	ExpiresAt  NullTime `json:"expires_at" dynamo:"expires_at"`
	CapturedAt NullTime `json:"captured_at" dynamo:"captured_at"`
	RefundedAt NullTime `json:"refunded_at" dynamo:"refunded_at"`
}
