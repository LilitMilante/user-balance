package entity

type Transfer struct {
	IDSender    int64 `json:"id_give"`
	IDRecipient int64 `json:"id_take"`
	Amount      int64 `json:"amount"`

	DescriptionSender    string `json:"-"`
	DescriptionRecipient string `json:"-"`
}
