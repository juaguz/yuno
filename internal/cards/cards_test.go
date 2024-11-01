package cards_test

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/google/uuid"
	"github.com/juaguz/yuno/internal/cards"
	"github.com/juaguz/yuno/internal/cards/dtos"
	"github.com/juaguz/yuno/internal/cards/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCardService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardRepo := mocks.NewMockCardRepository(ctrl)
	mockKmsRepo := mocks.NewMockKmsRepository(ctrl)
	mockVaultRepo := mocks.NewMockVaultRepository(ctrl)

	service := cards.NewCardService(mockCardRepo, mockKmsRepo, mockVaultRepo)

	userId := uuid.New()
	card := &dtos.Card{
		UserId: userId,
		Pan:    "encrypted_pan_data",
	}

	decryptedPan := base64.StdEncoding.EncodeToString([]byte("4111111111111111"))
	mockKmsRepo.EXPECT().Decrypt(gomock.Any(), "vault:v1:"+card.Pan, userId.String()).Return(decryptedPan, nil)
	mockCardRepo.EXPECT().Create(gomock.Any(), card).Return(nil)
	mockVaultRepo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	createdCard, err := service.Create(context.Background(), card)

	assert.NoError(t, err)
	assert.Equal(t, "4111", createdCard.Pan)
}

func TestCardService_Create_InvalidPAN(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardRepo := mocks.NewMockCardRepository(ctrl)
	mockKmsRepo := mocks.NewMockKmsRepository(ctrl)
	mockVaultRepo := mocks.NewMockVaultRepository(ctrl)

	service := cards.NewCardService(mockCardRepo, mockKmsRepo, mockVaultRepo)

	userId := uuid.New()
	card := &dtos.Card{
		UserId: userId,
		Pan:    "encrypted_pan_data",
	}

	decryptedPan := base64.StdEncoding.EncodeToString([]byte("1234567890123456"))
	mockKmsRepo.EXPECT().Decrypt(gomock.Any(), "vault:v1:"+card.Pan, userId.String()).Return(decryptedPan, nil)

	createdCard, err := service.Create(context.Background(), card)

	assert.ErrorIs(t, err, cards.ErrInvalidPan)
	assert.Nil(t, createdCard)
}

func TestCardService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardRepo := mocks.NewMockCardRepository(ctrl)
	mockKmsRepo := mocks.NewMockKmsRepository(ctrl)
	mockVaultRepo := mocks.NewMockVaultRepository(ctrl)

	service := cards.NewCardService(mockCardRepo, mockKmsRepo, mockVaultRepo)

	cardId := uuid.New()
	expectedCard := &dtos.Card{
		ID:     cardId,
		UserId: uuid.New(),
		Pan:    "4111",
	}

	mockCardRepo.EXPECT().Get(gomock.Any(), cardId).Return(expectedCard, nil)

	card, err := service.Get(context.Background(), expectedCard)

	assert.NoError(t, err)
	assert.Equal(t, expectedCard, card)
}

func TestCardService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardRepo := mocks.NewMockCardRepository(ctrl)
	mockKmsRepo := mocks.NewMockKmsRepository(ctrl)
	mockVaultRepo := mocks.NewMockVaultRepository(ctrl)

	service := cards.NewCardService(mockCardRepo, mockKmsRepo, mockVaultRepo)

	card := &dtos.Card{
		ID:     uuid.New(),
		UserId: uuid.New(),
		Pan:    "4111",
	}

	mockCardRepo.EXPECT().Get(gomock.Any(), card.ID).Return(card, nil)
	mockCardRepo.EXPECT().UpdateOne(gomock.Any(), card).Return(nil)

	err := service.Update(context.Background(), card)

	assert.NoError(t, err)
}

func TestCardService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardRepo := mocks.NewMockCardRepository(ctrl)
	mockKmsRepo := mocks.NewMockKmsRepository(ctrl)
	mockVaultRepo := mocks.NewMockVaultRepository(ctrl)

	service := cards.NewCardService(mockCardRepo, mockKmsRepo, mockVaultRepo)

	card := &dtos.Card{
		ID:     uuid.New(),
		UserId: uuid.New(),
		Pan:    "4111",
	}

	mockCardRepo.EXPECT().Get(gomock.Any(), card.ID).Return(card, nil)
	mockCardRepo.EXPECT().Delete(gomock.Any(), card.ID).Return(nil)
	mockVaultRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)

	err := service.Delete(context.Background(), card)

	assert.NoError(t, err)
}
