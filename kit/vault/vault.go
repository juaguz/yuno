package vault

import (
	"context"
	"fmt"

	vault "github.com/hashicorp/vault/api"
)

type VaultService struct {
	client   *vault.Client
	basePath string
}

func NewVaultService(client *vault.Client) *VaultService {
	return &VaultService{
		client:   client,
		basePath: "secret/data",
	}
}

func (v *VaultService) Create(ctx context.Context, data map[string]interface{}, key string) error {
	fullPath := fmt.Sprintf("%s/%s", v.basePath, key)

	secretData := map[string]interface{}{
		"data": data,
	}

	_, err := v.client.Logical().Write(fullPath, secretData)
	if err != nil {
		return fmt.Errorf("error al crear el secreto en Vault: %w", err)
	}

	return nil
}

func (v *VaultService) Delete(ctx context.Context, key string) error {
	fullPath := fmt.Sprintf("%s/%s", v.basePath, key)
	_, err := v.client.Logical().Delete(fullPath)
	if err != nil {
		return fmt.Errorf("error al eliminar el secreto en Vault: %w", err)
	}

	return nil
}
