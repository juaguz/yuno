package api

type CardCreation struct {
	CardHolder string `json:"card_holder"`
	Pan        string `json:"pan"`
}

type CardUpdate struct {
	CardHolder string `json:"card_holder"`
}
