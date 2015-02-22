package sshkeys

import "testing"

func TestCipher(t *testing.T) {
	pk := &PrivateKey{Path: "n/a", Block: aesRSAPEM}
	if pk.cipher() != "AES-128-CBC" {
		t.Error("expected AES-128-CBC, got", pk.cipher())
	}

	pk = &PrivateKey{Path: "n/a", Block: des3RSAPEM}
	if pk.cipher() != "DES-EDE3-CBC" {
		t.Error("expected DES-EDE3-CBC, got", pk.cipher())
	}
}

func TestDecryptAES(t *testing.T) {
	pk := &PrivateKey{Path: "n/a", Block: aesRSAPEM}
	pk.PasswordMethod = DummyPasswordMethod(passphrase)
	err := pk.Decrypt()
	if err != nil {
		t.Error("expect nil, got", err)
	}
}

func TestDecryptAESMacOS(t *testing.T) {
	pk := &PrivateKey{Path: "n/a", Block: aesRSAPEMMacOSMaverick}
	pk.PasswordMethod = DummyPasswordMethod(passphrase)
	err := pk.Decrypt()
	if err != nil {
		t.Error("expect nil, got", err)
	}

	pk = &PrivateKey{Path: "n/a", Block: aesRSAPEMMacOSYosemite}
	pk.PasswordMethod = DummyPasswordMethod(macOSYosemitePassphrase)
	err = pk.Decrypt()
	if err != nil {
		t.Error("expect nil, got", err)
	}
}

func TestDecryptDES3(t *testing.T) {
	pk := &PrivateKey{Path: "n/a", Block: des3RSAPEM}
	pk.PasswordMethod = DummyPasswordMethod(passphrase)
	err := pk.Decrypt()
	if err != nil {
		t.Error("expect nil, got", err)
	}

	pk = &PrivateKey{Path: "n/a", Block: des3RSAPEMLongPassphrase}
	pk.PasswordMethod = DummyPasswordMethod(longPassphrase)
	err = pk.Decrypt()
	if err != nil {
		t.Error("expect nil, got", err)
	}
}
