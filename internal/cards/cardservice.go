package cards

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/juaguz/yuno/internal/cards/dtos"
	"github.com/juaguz/yuno/kit/errors/senital"
)

var (
	ErrInvalidPan = errors.New("invalid pan")
)

func isValidCreditCard(cardNumber string) bool {
	regex := `^(?:4[0-9]{12}(?:[0-9]{3})?` + // Visa
		`|5[1-5][0-9]{14}` + // MasterCard
		`|3[47][0-9]{13}` + // American Express
		`|6(?:011|5[0-9]{2})[0-9]{12}` + // Discover
		`|(?:2131|1800|35\d{3})\d{11})$` // JCB

	re := regexp.MustCompile(regex)
	return re.MatchString(cardNumber)
}

func buildKey(userID uuid.UUID, cardID uuid.UUID) string {
	key := fmt.Sprintf("/secrets/cards/%s/%s", userID, cardID)
	return key
}

type CardRepository interface {
	Create(ctx context.Context, card *dtos.Card) error
	Get(ctx context.Context, id uuid.UUID) (*dtos.Card, error)
	UpdateOne(ctx context.Context, card *dtos.Card) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type KmsRepository interface {
	Decrypt(ctx context.Context, data string, key string) (string, error)
}

type VaultRepository interface {
	Create(ctx context.Context, data map[string]interface{}, key string) error
	Delete(ctx context.Context, key string) error
}

type CardService struct {
	CardRepository  CardRepository
	KmsRepository   KmsRepository
	VaultRepository VaultRepository
}

func NewCardService(cardRepository CardRepository, kmsRepository KmsRepository, vaultRepository VaultRepository) *CardService {
	return &CardService{
		CardRepository:  cardRepository,
		KmsRepository:   kmsRepository,
		VaultRepository: vaultRepository,
	}
}

func (c *CardService) Create(ctx context.Context, card *dtos.Card) (*dtos.Card, error) {
	card.ID = uuid.New()

	encryptedPan := card.Pan

	decryptedPan, err := c.KmsRepository.Decrypt(ctx, "vault:v1:"+encryptedPan, card.UserId.String())
	if err != nil {
		return nil, err
	}

	decodedPan, err := base64.StdEncoding.DecodeString(decryptedPan)
	if err != nil {
		return nil, err
	}

	pan := string(decodedPan)
	if !isValidCreditCard(pan) {
		return nil, ErrInvalidPan
	}

	card.Pan = pan[:4]

	err = c.CardRepository.Create(ctx, card)
	if err != nil {
		return nil, err
	}

	key := buildKey(card.UserId, card.ID)
	err = c.VaultRepository.Create(ctx, map[string]interface{}{
		"pan": encryptedPan,
	}, key)
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (c *CardService) Get(ctx context.Context, card *dtos.Card) (*dtos.Card, error) {
	card, err := c.CardRepository.Get(ctx, card.ID)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, senital.ErrNotFound
	}

	if card.UserId != card.UserId {
		return nil, senital.ErrNotFound
	}

	return card, nil
}

func (c *CardService) Update(ctx context.Context, card *dtos.Card) error {
	if _, err := c.Get(ctx, card); err != nil {
		return err
	}

	return c.CardRepository.UpdateOne(ctx, card)
}

func (c *CardService) Delete(ctx context.Context, card *dtos.Card) error {
	if _, err := c.Get(ctx, card); err != nil {
		return err
	}

	if err := c.CardRepository.Delete(ctx, card.ID); err != nil {
		return err
	}

	key := buildKey(card.UserId, card.ID)
	if err := c.VaultRepository.Delete(ctx, key); err != nil {
		return err
	}

	return nil
}
