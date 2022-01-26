package entity

type Balance struct {
	UserID int64 `json:"user_id"`
	Amount int64 `json:"amount"`
	TypeOp int   `json:"type_op,omitempty"`
}

const (
	Plus  = 1
	Minus = 2
)
