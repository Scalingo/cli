package db

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"

	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/term"
	"gopkg.in/errgo.v1"
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

	privateKey, err := sshPrivateKey(identity)
	if err != nil {
		return errgo.Mask(err)
	}

	sshConfig := &ssh.ClientConfig{
		User: "git",
		Auth: []ssh.AuthMethod{ssh.PublicKeys(privateKey)},
	}

	client, err := ssh.Dial("tcp", config.C.SshHost, sshConfig)
	if err != nil {
		return errgo.Mask(err)
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
		if env.Name == dbEnvVar {
			return env.Value
		}
	}
	return ""
}

func sshPrivateKey(path string) (ssh.Signer, error) {
	privateKey, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	// We parse the private key on our own first so that we can
	// show a nicer error if the private key has a password.
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, fmt.Errorf(
			"Failed to read key '%s': no key found", path)
	}
	if block.Headers["Proc-Type"] == "4,ENCRYPTED" {
		splitCipher := strings.Split(block.Headers["DEK-Info"], ",")
		cipherType, iv := splitCipher[0], strings.TrimSpace(splitCipher[1])
		switch cipherType {
		case "AES-128-CBC":
			decryptedKey, err := decryptAES128Key(block.Bytes, iv)
			if err != nil {
				return nil, errgo.Newf("Key is tagged AES-128-CBC, but is not: %v", err)
			}
			decryptedBlock := &pem.Block{}
			decryptedBlock.Type = block.Type
			decryptedBlock.Bytes = decryptedKey
			privateKey = pem.EncodeToMemory(decryptedBlock)
		default:
			return nil, fmt.Errorf(
				"Failed to read key '%s': password protected keys with '%s' are\n"+
					"not supported. Please decrypt the key prior to use.", path, cipherType)
		}
	}

	privateKeySigner, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, errgo.Newf("Invalid SSH key or password: %v", err)
	}

	return privateKeySigner, nil
}

func decryptAES128Key(payload []byte, iv string) ([]byte, error) {
	password, err := term.Password("Encrypted SSH Key, password: ")
	if err != nil {
		return nil, errgo.Mask(err)
	}

	ivBytes, err := hex.DecodeString(iv)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	key := append([]byte(password), ivBytes[0:8]...)
	keyHash := md5.New()
	_, err = keyHash.Write(key)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	keyHashed := keyHash.Sum(nil)

	decryptedPayload := make([]byte, len(payload))
	block, err := aes.NewCipher(keyHashed)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	decrypter := cipher.NewCBCDecrypter(block, ivBytes)
	decrypter.CryptBlocks(decryptedPayload, payload)
	decryptedPayload = bytes.TrimRight(decryptedPayload, "\x08\x0a")
	return decryptedPayload, nil
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
