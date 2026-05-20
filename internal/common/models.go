package common

import "github.com/google/uuid"

// PackSizeBatch is the request body for replacing all configured pack sizes.
type PackSizeBatch struct {
	Sizes []int `json:"sizes"`
}

// Order is the request and response shape for pack calculation (event_id is set on success).
type Order struct {
	ID          uuid.UUID `json:"event_id"`
	AmountItems int       `json:"items"`
}
