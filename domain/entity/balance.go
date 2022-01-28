package entity

type Balance struct {
	IDSender    int64  `json:"id_sender"`
	IDRecipient int64  `json:"id_recipient"`
	Amount      int64  `json:"amount"`
	TypeOp      int    `json:"type_op,omitempty"`
	Description string `json:"description,omitempty"`
}

const (
	Plus  = 1
	Minus = 2
)
