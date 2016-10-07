package sshkeys

import (
	"io"
	"net"
	"os"

	"github.com/Scalingo/cli/config"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"gopkg.in/errgo.v1"
)

func Client() (agent.Agent, net.Conn, error) {
	agentPath := os.Getenv("SSH_AUTH_SOCK")
	agentHandler, err := net.Dial("unix", agentPath)
	if err != nil {
		return nil, nil, errgo.Mask(err)
	}

	client := agent.NewClient(agentHandler)
	return client, agentHandler, nil
}

func ReadPrivateKeysFromAgent() ([]ssh.Signer, io.Closer, error) {
	client, agentHandler, err := Client()
	if err != nil {
		return nil, nil, errgo.Newf("Fail to communicate with SSH agent: %v\nPlease precise the SSH key you want to use with the flag -i", err)
	}
	config.C.Logger.Println("Using SSH agent to access private keys")

	signers, err := client.Signers()
	if err != nil {
		return nil, nil, errgo.Newf("Fail to access SSH keys throught SSH Agent: %v\n Please precise the SSH key you want to use with the flag -i", err)
	}
	return signers, agentHandler, nil
}

func AddKeyToAgent(privateKey interface{}) error {
	client, _, err := Client()
	if err == nil {
		key := agent.AddedKey{
			PrivateKey:       privateKey,
			ConfirmBeforeUse: false,
			LifetimeSecs:     0,
		}
		err := client.Add(key)
		if err != nil {
			return errgo.Mask(err)
		}
	}
	return nil
}
