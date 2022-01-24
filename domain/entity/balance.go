package entity

type Balance struct {
	UserID int `json:"user_id"`
	Amount int `json:"amount"`
	TypeOp int `json:"type_op"`
}

const (
	Plus  = 1
	Minus = 2
)
