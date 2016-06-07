package db

import (
	"errors"
	"fmt"
	stdio "io"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/cli/Godeps/_workspace/src/golang.org/x/crypto/ssh"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/debug"
	"github.com/Scalingo/cli/io"
	netssh "github.com/Scalingo/cli/net/ssh"
)

var (
	errTimeout      = errors.New("timeout")
	connIDGenerator = make(chan int)
	defaultPort     = 10000
)

type TunnelOpts struct {
	App       string
	DBEnvVar  string
	Identity  string
	Port      int
	Reconnect bool
}

func Tunnel(opts TunnelOpts) error {
	c := config.ScalingoClient()
	environ, err := c.VariablesListWithoutAlias(opts.App)
	if err != nil {
		return errgo.Mask(err)
	}

	dbUrlStr := dbEnvVarValue(opts.DBEnvVar, environ)
	if dbUrlStr == "" {
		return errgo.Newf("no such environment variable: %s", opts.DBEnvVar)
	}

	dbUrl, err := url.Parse(dbUrlStr)
	if err != nil {
		return errgo.Notef(err, "invalid database 'URL': %s", dbUrlStr)
	}
	fmt.Printf("Building tunnel to %s\n", dbUrl.Host)

	client, key, err := netssh.Connect(opts.Identity)
	if err != nil {
		return errgo.Notef(err, "fail to connect to SSH server")
	}
	waitingConnectionM := &sync.Mutex{}

	if opts.Port == 0 {
		opts.Port = defaultPort
	}

	var tcpAddr *net.TCPAddr
	var sock *net.TCPListener
	for {
		tcpAddr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", opts.Port))
		if err != nil {
			return errgo.Mask(err)
		}

		sock, err = net.ListenTCP("tcp", tcpAddr)
		if isAddrInUse(err) {
			opts.Port++
			continue
		}
		if err != nil {
			return errgo.Mask(err)
		}
		break
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

		debug.Println("Waiting local connection request")
		connToTunnel, err := sock.AcceptTCP()
		if err != nil {
			return errgo.Mask(err)
		}
		debug.Println("New local connection")
		// Checking not in reconnection process
		waitingConnectionM.Lock()
		waitingConnectionM.Unlock()

		go func() {
			for {
				err := handleConnToTunnel(client, dbUrl, connToTunnel, errs)
				if err != nil {
					debug.Println("Error happened in tunnel", err)
					if !opts.Reconnect {
						errs <- err
						return
					}
				}
				if err == errTimeout {
					waitingConnectionM.Lock()
					fmt.Println("Connection broken, reconnecting...")
					for err != nil {
						client, err = netssh.ConnectToSSHServerWithKey(key)
						if err != nil {
							fmt.Println("Fail to reconnect, waiting 10 seconds...")
							time.Sleep(10 * time.Second)
						}
					}
					fmt.Println("Reconnected!")
					waitingConnectionM.Unlock()
				}
				break
			}
		}()
	}
}

func dbEnvVarValue(dbEnvVar string, environ scalingo.Variables) string {
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

func handleConnToTunnel(sshClient *ssh.Client, dbUrl *url.URL, sock net.Conn, errs chan error) error {
	connID := <-connIDGenerator
	fmt.Printf("Connect to %s [%v]\n", dbUrl.Host, connID)
	conn, err := sshClient.Dial("tcp", dbUrl.Host)
	if err != nil {
		errs <- err
		return nil
	}
	debug.Println("Connected to", dbUrl.Host, connID)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		debug.Println("Pipe DB -> Local ON")
		_, remoteErr := stdio.Copy(sock, conn)
		debug.Println("Pipe DB -> Local OFF", remoteErr)
		sock.Close()
		wg.Done()
	}()

	go func() {
		debug.Println("Local -> DB ON")
		_, err = io.CopyWithTimeout(2*time.Second)(conn, sock)
		debug.Println("Local -> DB OFF", err)
		conn.Close()
		wg.Done()
	}()

	wg.Wait()

	fmt.Printf("End of connection [%d]\n", connID)
	// Connection timeout
	if err != nil && strings.Contains(err.Error(), "use of closed network") {
		return errTimeout
	}
	return nil
}

func startIDGenerator() {
	for i := 1; ; i++ {
		connIDGenerator <- i
	}
}

func isAddrInUse(err error) bool {
	if err, ok := err.(*net.OpError); ok {
		if err, ok := err.Err.(*os.SyscallError); ok {
			return err.Err == syscall.EADDRINUSE
		}
	}
	return false
}
