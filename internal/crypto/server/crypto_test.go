package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRSAPrivateKey(t *testing.T) {
	t.Run("should return RSA private key when file exists and contains valid PEM block with RSA private key", func(t *testing.T) {
		validPEM, err := generateValidRSAPrivateKeyPEM()
		require.NoError(t, err, "unexpected error when generate valid private key")
		tmp, err := os.CreateTemp("", "private-key-*.pem")
		require.NoError(t, err, "unexpected error when create temp file")
		_, err = tmp.Write([]byte(validPEM))
		require.NoError(t, err, "unexpected error when write private key to temp file")
		err = tmp.Close()
		require.NoError(t, err, "unexpected error when close temp file")
		defer os.Remove(tmp.Name())

		key, err := GetRSAPrivateKey(tmp.Name())
		assert.NoError(t, err, "unexpected error when get rsa private key")
		assert.NotNil(t, key, "rsa private key should not be nil")
	})

	t.Run("should return error when file not exists", func(t *testing.T) {
		_, err := GetRSAPrivateKey("not_exists_file.pem")

		assert.Error(t, err, "should return error")
	})

	t.Run("should return error when file exists and contains invalid PEM block with RSA private key", func(t *testing.T) {
		invalidPEM := generateInvalidPEM()
		tmp, err := os.CreateTemp("", "private-key-*.pem")
		require.NoError(t, err, "unexpected error when create temp file")
		_, err = tmp.Write([]byte(invalidPEM))
		require.NoError(t, err, "unexpected error when write private key to temp file")
		err = tmp.Close()
		require.NoError(t, err, "unexpected error when close temp file")
		defer os.Remove(tmp.Name())

		_, err = GetRSAPrivateKey(tmp.Name())

		assert.Error(t, err, "should return error")
	})
}

func generateValidRSAPrivateKeyPEM() (string, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}
	privateKey := x509.MarshalPKCS1PrivateKey(key)
	pemBlock := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKey,
	})
	return string(pemBlock), nil
}

func generateInvalidPEM() string {
	return `-----BEGIN RSA PRIVATE KEY-----
invalid PEM block
-----END RSA PRIVATE KEY-----`
}
