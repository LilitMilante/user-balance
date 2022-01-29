package entity

import "time"

type Transaction struct {
	ID         int64     `json:"id,omitempty"`
	UserID     int64     `json:"user_id"`
	Amount     int64     `json:"amount"`
	Event      int       `json:"event"`
	TransferID int64     `json:"transfer_id"`
	Message    string    `json:"message,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}

const (
	EventCrediting = 1
	EventWriteOffs = 2
	EventTransfer  = 3
)
