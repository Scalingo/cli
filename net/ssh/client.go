package ssh

import (
	"context"
	stdio "io"

	"golang.org/x/crypto/ssh"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/crypto/sshkeys"
	"github.com/Scalingo/go-scalingo/v8/debug"
)

var (
	ErrNoAuthSucceed = errgo.Newf("No authentication method has succeeded")
)

type ConnectOpts struct {
	Host     string
	Identity string
}

func Connect(ctx context.Context, opts ConnectOpts) (*ssh.Client, ssh.Signer, error) {
	var (
		err         error
		privateKeys []ssh.Signer
	)
	if opts.Identity == "ssh-agent" {
		var agentConnection stdio.Closer
		privateKeys, agentConnection, err = sshkeys.ReadPrivateKeysFromAgent()
		if err != nil {
			return nil, nil, errgo.Mask(err)
		}
		defer agentConnection.Close()
	}

	if len(privateKeys) == 0 {
		if opts.Identity == "ssh-agent" {
			opts.Identity = sshkeys.DefaultKeyPath
		}
		privateKey, err := sshkeys.ReadPrivateKey(ctx, opts.Identity)
		if err != nil {
			return nil, nil, errgo.Mask(err)
		}
		privateKeys = append(privateKeys, privateKey)
	}

	debug.Println("Identity used:", opts.Identity, "Private keys:", len(privateKeys))

	client, key, err := connectToSSHServer(connectSSHOpts{
		Host: opts.Host,
		Keys: privateKeys,
	})
	if err != nil {
		return nil, nil, err
	}
	debug.Println("SSH connection:", client.LocalAddr(), "Key:", string(key.PublicKey().Marshal()))
	return client, key, nil
}

type connectSSHOpts struct {
	Host string
	Keys []ssh.Signer
}

func connectToSSHServer(opts connectSSHOpts) (*ssh.Client, ssh.Signer, error) {
	var (
		client     *ssh.Client
		privateKey ssh.Signer
		err        error
	)

	for _, privateKey = range opts.Keys {
		client, err = connectToSSHServerWithKey(opts.Host, privateKey)
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

func connectToSSHServerWithKey(host string, key ssh.Signer) (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User:            "git",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(key)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return ssh.Dial("tcp", host, sshConfig)
}
