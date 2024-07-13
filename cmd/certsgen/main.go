// Package is generator of RSA keys.
//
// It generates an RSA private-public key pair, encodes them in PEM format,
// and writes them to files named "private.pem" and "public.pem" in a directory "./certs".
//
// Example usage:
//
//	$ go run .
//
// The generated files can then be used for cryptographic operations such as signing
// and verification, encryption and decryption.
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	err = os.WriteFile("./certs/private.pem", privateKeyBlock, 0644)
	if err != nil {
		log.Fatal(err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Fatal(err)
	}
	publicKeyBlock := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	err = os.WriteFile("./certs/public.pem", publicKeyBlock, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
