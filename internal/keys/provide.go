package keys

import (
	"context"

	"github.com/google/uuid"
)

type KmsRepo interface {
	GetPublicKey(ctx context.Context, keyID string) (string, error)
	CreateKey(ctx context.Context, keyID string) error
}

type KeysProvider struct {
	KmsRepo KmsRepo
}

func NewKeysProvider(kmsRepo KmsRepo) *KeysProvider {
	return &KeysProvider{KmsRepo: kmsRepo}
}

func (k *KeysProvider) CreateKey(ctx context.Context, userID uuid.UUID) (string, error) {
	keyID := userID.String()
	err := k.KmsRepo.CreateKey(ctx, keyID)
	if err != nil {
		return "", err
	}

	return k.KmsRepo.GetPublicKey(ctx, keyID)
}
