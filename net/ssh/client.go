package ssh

import (
	stdio "io"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/crypto/sshkeys"
	"github.com/Scalingo/cli/debug"
	"golang.org/x/crypto/ssh"
	"gopkg.in/errgo.v1"
)

var (
	ErrNoAuthSucceed = errgo.Newf("No authentication method has succeeded")
)

func Connect(identity string) (*ssh.Client, ssh.Signer, error) {
	var (
		err         error
		privateKeys []ssh.Signer
	)
	if identity == "ssh-agent" {
		var agentConnection stdio.Closer
		privateKeys, agentConnection, err = sshkeys.ReadPrivateKeysFromAgent()
		if err != nil {
			return nil, nil, errgo.Mask(err)
		}
		defer agentConnection.Close()
	}

	if len(privateKeys) == 0 {
		if identity == "ssh-agent" {
			identity = sshkeys.DefaultKeyPath
		}
		privateKey, err := sshkeys.ReadPrivateKey(identity)
		if err != nil {
			return nil, nil, errgo.Mask(err)
		}
		privateKeys = append(privateKeys, privateKey)
	}

	debug.Println("Identity used:", identity, "Private keys:", len(privateKeys))

	client, key, err := ConnectToSSHServer(privateKeys)
	if err != nil {
		return nil, nil, err
	}
	debug.Println("SSH connection:", client.LocalAddr(), "Key:", string(key.PublicKey().Marshal()))
	return client, key, nil
}

func ConnectToSSHServer(keys []ssh.Signer) (*ssh.Client, ssh.Signer, error) {
	var (
		client     *ssh.Client
		privateKey ssh.Signer
		err        error
	)

	for _, privateKey = range keys {
		client, err = ConnectToSSHServerWithKey(privateKey)
		if err == nil {
			break
		} else {
			config.C.Logger.Println("Fail to connect to the SSH server", err)
		}
	}
	if client == nil {
		return nil, nil, ErrNoAuthSucceed
	}
	return client, privateKey, nil
}

func ConnectToSSHServerWithKey(key ssh.Signer) (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: "git",
		Auth: []ssh.AuthMethod{ssh.PublicKeys(key)},
	}

	return ssh.Dial("tcp", config.C.SshHost, sshConfig)
}
