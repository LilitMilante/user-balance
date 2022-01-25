package entity

type Balance struct {
	UserID int     `json:"user_id"`
	Amount float64 `json:"amount"`
	TypeOp int     `json:"type_op,omitempty"`
}

const (
	Plus  = 1
	Minus = 2
)
