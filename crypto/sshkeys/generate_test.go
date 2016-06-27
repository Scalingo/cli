package sshkeys

import (
	"crypto/rand"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestGeneration(t *testing.T) {
	publicKeyByte, privateKeyByte, err := GenerateKey()
	if err != nil {
		t.Error("expected nil, got", err)
	}

	publicKey, _, _, _, err := ssh.ParseAuthorizedKey(publicKeyByte)
	if err != nil {
		t.Error("expected nil, got", err)
	}

	privateKeySigner, err := ssh.ParsePrivateKey(privateKeyByte)
	if err != nil {
		t.Error("expected nil, got", err)
	}

	payload := []byte("Hello World")
	sign, err := privateKeySigner.Sign(rand.Reader, payload)
	if err != nil {
		t.Error("expected nil, got", err)
	}

	err = publicKey.Verify(payload, sign)

	if err != nil {
		t.Error("expected nil, got", err)
	}
}
