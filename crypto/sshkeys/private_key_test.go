package sshkeys

import (
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/require"
)

func pemDecode(keyToDecode []byte) *pem.Block {
	block, _ := pem.Decode(keyToDecode)
	return block
}

func TestPrivateKey_IsEncrypted(t *testing.T) {
	cases := []struct {
		Name          string
		PrivateKey    *PrivateKey
		ExpectBoolean bool
	}{
		{
			Name:          "Unencrypted RSA Key",
			PrivateKey:    &PrivateKey{Path: "n/a", Block: pemDecode(unencryptedRSAKey)},
			ExpectBoolean: false,
		},
		{
			Name:          "Encrypted AES RSA Key",
			PrivateKey:    &PrivateKey{Path: "n/a", Block: pemDecode(aesRSA)},
			ExpectBoolean: true,
		},
		{
			Name:          "Encrypted DES3 RSA Key",
			PrivateKey:    &PrivateKey{Path: "n/a", Block: pemDecode(des3RSA)},
			ExpectBoolean: true,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			boolean := c.PrivateKey.IsEncrypted()

			require.Equal(t, c.ExpectBoolean, boolean)
		})
	}
}
