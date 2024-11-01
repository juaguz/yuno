package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
)

func main() {
	// Clave pública en formato PEM
	publicKeyPEM := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAq/GbunxrZxT6SJlV8eDH
kbA9efN387baQv6AMXxknDXWtLD0bNpwlwFEh/fbL17FZu6EvhipZiGegTf1OvHF
cVt0xiEE7oJ85YM66AW/U4GpksqcTGWbenb24IEzN0iHs7+f+BFQtYYu4vo2mTjT
6/olc2YFajEdIGyohYtqZ38bwrAeSw+bxCDfsbDxb8OWVuncsBO4P1siQ4GhpDTv
Mplk1AvrWtr4BTu75yDJHqkUr6RyoA0MWp1nBLC5b3KXegS2djMGqRTTAl7SQdfU
3dGJEL3Bu4JFXerGfdwNZ8YRtdW5+jpgsxhaaKO03q74YhyRxwFvpSSi89DWZxex
HwIDAQAB
-----END PUBLIC KEY-----`

	// Texto a cifrar (ejemplo de PAN)
	plainText := "4140313934452647"

	// Codificar el texto en Base64
	plainTextBase64 := plainText

	// Parsear la clave pública desde el formato PEM
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Fatalf("Error al cargar la clave pública: formato no compatible")
	}

	// Convertir la clave pública PEM en una clave RSA
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("Error al parsear la clave pública: %v", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		log.Fatalf("Clave pública no es de tipo RSA")
	}

	// Cifrar el texto en Base64 utilizando RSA-OAEP con SHA-256
	label := []byte("")
	hash := sha256.New()
	encryptedBytes, err := rsa.EncryptOAEP(hash, rand.Reader, rsaPublicKey, []byte(plainTextBase64), label)
	if err != nil {
		log.Fatalf("Error al cifrar el texto: %v", err)
	}

	// Codificar el cifrado en Base64 para transportarlo o almacenarlo
	encryptedTextBase64 := base64.StdEncoding.EncodeToString(encryptedBytes)

	fmt.Println("Texto en claro (Base64):", plainTextBase64)
	fmt.Println("Texto cifrado (Base64):", encryptedTextBase64)
}
