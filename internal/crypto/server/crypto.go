package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// GetRSAPrivateKey reads an RSA private key from a PEM-encoded file.
//
// This function takes the file path of a PEM-encoded RSA private key as an argument,
// reads the file, decodes the PEM block, and parses the private key. It returns a pointer
// to the parsed rsa.PrivateKey and an error, if any occurs during the process.
func GetRSAPrivateKey(fPath string) (*rsa.PrivateKey, error) {
	privateKeyPEM, err := os.ReadFile(fPath)
	if err != nil {
		return nil, fmt.Errorf("read file with private key for asymmetric encryption: %w", err)
	}

	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key for asymmetric encryption: %w", err)
	}

	return rsaPrivateKey, nil
}
