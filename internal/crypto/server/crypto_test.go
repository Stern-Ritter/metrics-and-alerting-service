package server

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRSAPrivateKey(t *testing.T) {
	t.Run("should return RSA private key when file exists and contains valid PEM block with RSA private key", func(t *testing.T) {
		validPEM := `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`
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
		invalidPEM := `-----BEGIN RSA PRIVATE KEY-----
invalid PEM block
-----END RSA PRIVATE KEY-----`
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
