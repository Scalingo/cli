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

func ReadPrivateKeysFromAgent() ([]ssh.Signer, io.Closer, error) {
	agentPath := os.Getenv("SSH_AUTH_SOCK")
	agentHandler, err := net.Dial("unix", agentPath)
	if err != nil {
		return nil, nil, errgo.Newf("Fail to communicate with SSH agent: %v\nPlease precise the SSH key you want to use with the flag -i", err)
	}
	config.C.Logger.Println("Using SSH agent to access private keys")

	client := agent.NewClient(agentHandler)
	signers, err := client.Signers()
	if err != nil {
		return nil, nil, errgo.Newf("Fail to access SSH keys through SSH Agent: %v\n Please precise the SSH key you want to use with the flag -i", err)
	}
	return signers, agentHandler, nil
}
