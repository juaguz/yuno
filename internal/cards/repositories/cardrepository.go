package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/juaguz/yuno/internal/cards/dtos"
	"github.com/juaguz/yuno/internal/cards/models"
	"github.com/juaguz/yuno/kit/database"
	"gorm.io/gorm"
)

type CardRepository struct {
	DB *gorm.DB
}

func NewCardRepository(DB *gorm.DB) *CardRepository {
	return &CardRepository{DB: DB}
}

func (c CardRepository) Create(ctx context.Context, card *dtos.Card) error {
	db := database.GetTx(ctx, c.DB)
	cc := &models.Card{
		CardHolder: card.CardHolder,
		UserId:     card.UserId,
		LastDigits: card.Pan,
	}

	return db.Create(cc).Error
}

func (c CardRepository) Get(ctx context.Context, id uuid.UUID) (*dtos.Card, error) {
	var cardModel models.Card
	if err := c.DB.First(&cardModel, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &dtos.Card{
		ID:         cardModel.ID,
		CardHolder: cardModel.CardHolder,
		Pan:        cardModel.LastDigits,
		UserId:     cardModel.UserId,
	}, nil
}

func (c CardRepository) UpdateOne(ctx context.Context, card *dtos.Card) error {
	db := database.GetTx(ctx, c.DB)
	updateData := models.Card{
		CardHolder: card.CardHolder,
	}

	if err := db.Model(&models.Card{}).Where("id = ?", card.ID).Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}

func (c CardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db := database.GetTx(ctx, c.DB)
	if err := db.Delete(&models.Card{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}
