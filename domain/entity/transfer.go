package entity

type Transfer struct {
	IdGive int64 `json:"id_give"`
	IdTake int64 `json:"id_take"`
	Amount int64 `json:"amount"`

	DescriptionSender    string `json:"-"`
	DescriptionRecipient string `json:"-"`
}
