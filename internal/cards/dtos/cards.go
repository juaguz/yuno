package dtos

import "github.com/google/uuid"

type Card struct {
	ID         uuid.UUID `json:"id"`
	CardHolder string    `json:"card_holder"`
	Pan        string    `json:"pan"`
	UserId     uuid.UUID `json:"user_id"`
}

//enum for status

type Status string

const (
	Succeeded Status = "succeeded"
	Failed    Status = "failed"
)

type BatchUpdateStatus struct {
	Status Status    `json:"status"`
	CardID uuid.UUID `json:"card"`
}

type BatchUpdate struct {
	ID         uuid.UUID `json:"id"`
	CardHolder string    `json:"card_holder"`
}
