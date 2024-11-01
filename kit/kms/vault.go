package kms

import (
	"context"
	"encoding/json"
	"fmt"

	vault "github.com/hashicorp/vault/api"
)

type VaultKmsService struct {
	client *vault.Client
}

func NewVaultKmsService(client *vault.Client) *VaultKmsService {
	return &VaultKmsService{client: client}
}

func (v *VaultKmsService) Decrypt(ctx context.Context, encryptedData, keyID string) (string, error) {
	transitPath := fmt.Sprintf("transit/decrypt/%s", keyID)

	data := map[string]interface{}{
		"ciphertext": encryptedData,
	}

	secret, err := v.client.Logical().Write(transitPath, data)
	if err != nil {
		return "", fmt.Errorf("error desencriptando datos: %w", err)
	}

	decryptedData, ok := secret.Data["plaintext"].(string)
	if !ok {
		return "", fmt.Errorf("error al obtener datos desencriptados desde Vault")
	}

	return decryptedData, nil
}

func (v *VaultKmsService) CreateKey(ctx context.Context, keyID string) error {
	transitPath := fmt.Sprintf("transit/keys/%s", keyID)

	data := map[string]interface{}{
		"type":                   "rsa-2048",
		"exportable":             true,
		"allow_plaintext_backup": false,
	}

	_, err := v.client.Logical().Write(transitPath, data)
	if err != nil {
		return fmt.Errorf("creating keys: %w", err)
	}

	return nil
}

func (v *VaultKmsService) GetPublicKey(ctx context.Context, keyID string) (string, error) {
	transitPath := fmt.Sprintf("transit/keys/%s", keyID)

	secret, err := v.client.Logical().Read(transitPath)
	if err != nil {
		return "", fmt.Errorf("getting public key: %w", err)
	}

	latestVersion, ok := secret.Data["latest_version"].(json.Number)
	if !ok {
		return "", fmt.Errorf("can't get latest version from Vault")
	}

	keys, ok := secret.Data["keys"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("can't get keys map from Vault")
	}

	versionKey := fmt.Sprintf("%s", latestVersion.String())
	keyData, ok := keys[versionKey].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("can't get key data for latest version from Vault")
	}

	publicKey, ok := keyData["public_key"].(string)
	if !ok {
		return "", fmt.Errorf("can't get public key from Vault")
	}

	return publicKey, nil
}
