package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

var (
	EcdsaPrivateKey *ecdsa.PrivateKey
	EcdsaPublicKey  *ecdsa.PublicKey
)

func loadECDSAPrivateKey(filePath string) (*ecdsa.PrivateKey, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to parse private key")
	}
	var privateKey *ecdsa.PrivateKey
	privateKey, err = x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func loadECDSAPublicKey(filePath string) (*ecdsa.PublicKey, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to parse public key")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := pubInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not an ECDSA public key")
	}

	return publicKey, nil
}

func InitEcdsa() {
	var err error
	EcdsaPrivateKey, err = loadECDSAPrivateKey("/root/project/backend/auth/ecdsa_private.pem")
	if err != nil {
		panic(err)
	}
	EcdsaPublicKey, err = loadECDSAPublicKey("/root/project/backend/auth/ecdsa_public.pem")
	if err != nil {
		panic(err)
	}
}
