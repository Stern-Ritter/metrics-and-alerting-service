package agent

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

func TestGetRSAPublicKey(t *testing.T) {
	t.Run("should return RSA public key when file exits and contains valid PEM block with RSA public key",
		func(t *testing.T) {
			validPEM, err := generateValidRSAPublicKeyPEM()
			require.NoError(t, err, "unexpect error when generate valid RSA public key")
			tmp, err := os.CreateTemp("", "public-key-*.pem")
			require.NoError(t, err, "unexpected error when create temp file")
			_, err = tmp.Write([]byte(validPEM))
			require.NoError(t, err, "unexpected error when write public key to temp file")
			err = tmp.Close()
			require.NoError(t, err, "unexpected error when close temp file")
			defer os.Remove(tmp.Name())

			key, err := GetRSAPublicKey(tmp.Name())

			assert.NoError(t, err, "unexpected error when get rsa public key")
			assert.NotNil(t, key, "rsa public key should not be nil")
		})

	t.Run("should return error when file not exists", func(t *testing.T) {
		_, err := GetRSAPublicKey("not_exists_file.pem")

		assert.Error(t, err, "should return error")
	})

	t.Run("should return error when file exits and contains invalid PEM block with RSA public key",
		func(t *testing.T) {
			invalidPEM := generateInvalidPEM()
			tmp, err := os.CreateTemp("", "public-key-*.pem")
			require.NoError(t, err, "unexpected error when create temp file")
			_, err = tmp.Write([]byte(invalidPEM))
			require.NoError(t, err, "unexpected error when write public key to temp file")
			err = tmp.Close()
			require.NoError(t, err, "unexpected error when close temp file")
			defer os.Remove(tmp.Name())

			_, err = GetRSAPublicKey(tmp.Name())

			assert.Error(t, err, "should return error")
		})

	t.Run("should return error when file exits and contains valid PEM block with non RSA public key", func(t *testing.T) {
		nonRSAPEM := generateNonRSAPublicKeyPEM()
		temp, err := os.CreateTemp("", "public-key-*.pem")
		require.NoError(t, err, "unexpected error when create temp file")
		_, err = temp.Write([]byte(nonRSAPEM))
		require.NoError(t, err, "unexpected error when write public key to temp file")
		err = temp.Close()
		require.NoError(t, err, "unexpected error when close temp file")
		defer os.Remove(temp.Name())

		_, err = GetRSAPublicKey(temp.Name())

		assert.Error(t, err, "should return error")
	})
}

func generateValidRSAPublicKeyPEM() (string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", err
	}
	key := &privateKey.PublicKey
	publicKey, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}
	pemBlock := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKey,
	})
	return string(pemBlock), nil
}

func generateInvalidPEM() string {
	return `-----BEGIN PUBLIC KEY-----
invalid PEM block
-----END PUBLIC KEY-----`
}

func generateNonRSAPublicKeyPEM() string {
	return `-----BEGIN PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEdH2FQczFvIQaiBgAcnlKSnmFsEAKrmFz
eB8yBZBsIr3jNkGeITIkDkAAQ0YyJGBB4jDqOG5r+3lKx7R8GhfTi/xt6n8LVZUS
hwlH3TfglRAVJg1dQxxXmnQXsp45Im1t
-----END PUBLIC KEY-----`
}
