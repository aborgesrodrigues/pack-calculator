package common

import "github.com/google/uuid"

type PackSizeBatch struct {
	Sizes []int `json:"sizes"`
}

type Order struct {
	ID          uuid.UUID `json:"event_id"`
	AmountItems int       `json:"items"`
}
