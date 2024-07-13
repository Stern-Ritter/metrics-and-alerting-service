package agent

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// GetRSAPublicKey reads an RSA public key from a PEM-encoded file.
//
// This function takes the file path of a PEM-encoded RSA public key as an argument,
// reads the file, decodes the PEM block, and parses the public key. It returns a pointer
// to the parsed rsa.PublicKey and an error, if any occurs during the process.
func GetRSAPublicKey(fPath string) (*rsa.PublicKey, error) {
	publicKeyPEM, err := os.ReadFile(fPath)
	if err != nil {
		return nil, fmt.Errorf("read file with public key for asymmetric encryption: %w", err)
	}

	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key for asymmetric encryption: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key for asymmetric encryption is not RSA")
	}

	return rsaPublicKey, nil
}
