package db

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"sync"

	"github.com/Scalingo/cli/Godeps/_workspace/src/golang.org/x/crypto/ssh"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/crypto/sshkeys"
	"github.com/Scalingo/cli/debug"
)

var (
	connIDGenerator = make(chan int)
)

func Tunnel(app string, dbEnvVar string, identity string, port int) error {
	environ, err := api.VariablesList(app)
	if err != nil {
		return errgo.Mask(err)
	}

	dbUrlStr := dbEnvVarValue(dbEnvVar, environ)
	if dbUrlStr == "" {
		return errgo.Newf("no such environment variable: %s", dbEnvVar)
	}

	dbUrl, err := url.Parse(dbUrlStr)
	if err != nil {
		return errgo.Notef(err, "invalid database 'URL': %s", dbUrlStr)
	}
	fmt.Printf("Building tunnel to %s\n", dbUrl.Host)

	var privateKeys []ssh.Signer
	if identity == "ssh-agent" {
		var agentConnection io.Closer
		privateKeys, agentConnection, err = sshkeys.ReadPrivateKeysFromAgent()
		if err != nil {
			return errgo.Mask(err)
		}
		defer agentConnection.Close()
	} else {
		privateKey, err := sshkeys.ReadPrivateKey(identity)
		if err != nil {
			return errgo.Mask(err)
		}
		debug.Println("Identity used:", identity)
		privateKeys = append(privateKeys, privateKey)
	}

	var client *ssh.Client
	for _, privateKey := range privateKeys {
		sshConfig := &ssh.ClientConfig{
			User: "git",
			Auth: []ssh.AuthMethod{ssh.PublicKeys(privateKey)},
		}

		client, err = ssh.Dial("tcp", config.C.SshHost, sshConfig)
		if err == nil {
			break
		} else {
			config.C.Logger.Println("Fail to connect to the SSH server", err)
		}
	}
	if client == nil {
		return errgo.Newf("Unable to connect to our SSH server: %v", err)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return errgo.Mask(err)
	}

	sock, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return errgo.Mask(err)
	}
	defer sock.Close()
	fmt.Printf("You can access your database on '%v'\n", sock.Addr())

	go startIDGenerator()
	errs := make(chan error)
	for {
		select {
		case err := <-errs:
			return errgo.Mask(err)
		default:
		}

		connToTunnel, err := sock.Accept()
		if err != nil {
			return errgo.Mask(err)
		}
		go handleConnToTunnel(client, dbUrl, connToTunnel, errs)
	}
}

func dbEnvVarValue(dbEnvVar string, environ api.Variables) string {
	for _, env := range environ {
		if env.Value == dbEnvVar {
			return dbEnvVar
		}
		if env.Name == dbEnvVar {
			return env.Value
		}
	}
	return ""
}

func handleConnToTunnel(sshClient *ssh.Client, dbUrl *url.URL, sock net.Conn, errs chan error) {
	connID := <-connIDGenerator
	fmt.Printf("Connect to %s [%v]\n", dbUrl.Host, connID)
	conn, err := sshClient.Dial("tcp", dbUrl.Host)
	if err != nil {
		errs <- err
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		io.Copy(sock, conn)
		sock.Close()
		wg.Done()
	}()

	go func() {
		io.Copy(conn, sock)
		conn.Close()
		wg.Done()
	}()

	wg.Wait()

	fmt.Printf("End of connection [%d]\n", connID)
}

func startIDGenerator() {
	for i := 1; ; i++ {
		connIDGenerator <- i
	}
}
