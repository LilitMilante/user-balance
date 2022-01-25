package entity

type Transfer struct {
	IdGive int `json:"id_give"`
	IdTake int `json:"id_take"`
	Amount int `json:"amount"`
}
