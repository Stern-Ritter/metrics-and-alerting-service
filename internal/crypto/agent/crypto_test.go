package agent

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRSAPublicKey(t *testing.T) {
	t.Run("should return RSA public key when file exits and contains valid PEM block with RSA public key",
		func(t *testing.T) {
			validPEM := `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEArB/xf+O7cxZDgl1+ndk2
m/tn3UotjmxQMj6k48uZR6THMEnO3XUB0cyUPXU80Odd/K1nxHyZwRRNDfTGeO61
yNsX7a6j+R4vnM2+qhyCQ/ymYqNZ6D2NyzqJl1XODL9DiFBoZImP12Y6j+Q4NkYB
2KY7dhDpliqfC/wtZ3FTczeHBetGy1A54mn351FV7+VcSrOuhu+KCz/fLBxeCcqq
oWgeFgLNPJnEwjejK+gX84O/Etak+dILldlcMQyiOapnSUHgIfRtPQrMjL4DtsrD
AQQ1+viz4g1OrF0KAEFC/yjnEoeX22rkm9oQO6RjnM1MqzP9sNTKin5N+mz9AV09
l/SYLczXaLggbt/ufLM5h6pDEmEH0gz616VK425fKP9CNaX9mlnIq8yvfcG1CL0t
kaoIhPDjZsi6MX/+mDYuNzJLqxQgtDXoQ6nAYfGkSEfYNNHIVz5REPXfn3TRtUia
gHx6Yhpj3988PD9JvAUrxKi/KHhjnWSNBJOQ8tEOrjTyPMYHi5pMdwFso/fzTHFt
MLwrVyx7EX+Pp9NKUNbnY9PfHPj2xnUUxGEWKSMbHjviGBfxoHewuPluN45m/7mR
/Mf+cWN8DWRyzN04cT00Ba7hI3NSubemVdM3wuBgEJ680VRx2IqEAo1cNc7OjYbT
mas6/wc1G8V6aMF0k06PSQMCAwEAAQ==
-----END PUBLIC KEY-----`
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
			invalidPEM := `-----BEGIN PUBLIC KEY-----
invalid PEM block
-----END PUBLIC KEY-----`
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
		nonRSAPEM := `-----BEGIN PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEdH2FQczFvIQaiBgAcnlKSnmFsEAKrmFz
eB8yBZBsIr3jNkGeITIkDkAAQ0YyJGBB4jDqOG5r+3lKx7R8GhfTi/xt6n8LVZUS
hwlH3TfglRAVJg1dQxxXmnQXsp45Im1t
-----END PUBLIC KEY-----`
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
