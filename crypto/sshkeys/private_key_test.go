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
			Name:          "Encrypted AES RSA Key (Mac OS Maverick)",
			PrivateKey:    &PrivateKey{Path: "n/a", Block: pemDecode(aesRSAMacOSMaverick)},
			ExpectBoolean: true,
		},
		{
			Name:          "Encrypted AES RSA Key (Mac OS Yosemite)",
			PrivateKey:    &PrivateKey{Path: "n/a", Block: pemDecode(aesRSAMacOSYosemite)},
			ExpectBoolean: true,
		},
		{
			Name:          "Encrypted DES3 RSA Key",
			PrivateKey:    &PrivateKey{Path: "n/a", Block: pemDecode(des3RSA)},
			ExpectBoolean: true,
		},
		{
			Name:          "Encrypted DES3 RSA Key with Long Passphrase",
			PrivateKey:    &PrivateKey{Path: "n/a", Block: pemDecode(des3RSALongPassphrase)},
			ExpectBoolean: true,
		},
		{
			Name:          "Unencrypted OpenSSH RSA Key",
			PrivateKey:    &PrivateKey{Path: "n/a", Block: pemDecode(unencryptedOpenSSHRSA)},
			ExpectBoolean: false,
		},
		{
			Name:          "Encrypted OpenSSH RSA Key",
			PrivateKey:    &PrivateKey{Path: "n/a", Block: pemDecode(openSSHRSA)},
			ExpectBoolean: true,
		},
		{
			Name:          "Unencrypted OpenSSH ecdsa Key",
			PrivateKey:    &PrivateKey{Path: "n/a", Block: pemDecode(unencryptedOpenSSHecdsa)},
			ExpectBoolean: false,
		},
		{
			Name:          "Encrypted OpenSSH ecdsa Key",
			PrivateKey:    &PrivateKey{Path: "n/a", Block: pemDecode(openSSHecdsa)},
			ExpectBoolean: true,
		},
		{
			Name:          "Unencrypted OpenSSH ed25519 Key",
			PrivateKey:    &PrivateKey{Path: "n/a", Block: pemDecode(unencryptedOpenSSHed25519)},
			ExpectBoolean: false,
		},
		{
			Name:          "Encrypted OpenSSH ed25519 Key",
			PrivateKey:    &PrivateKey{Path: "n/a", Block: pemDecode(openSSHed25519)},
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
