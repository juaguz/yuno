package models

import (
	"github.com/google/uuid"
	"github.com/juaguz/yuno/kit/database"
)

type Card struct {
	database.Model
	CardHolder string    `json:"card_holder"`
	UserId     uuid.UUID `json:"user_id"`
	LastDigits string    `json:"last_digits"`
}
