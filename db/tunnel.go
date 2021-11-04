package db

import (
	"fmt"
	stdio "io"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	netssh "github.com/Scalingo/cli/net/ssh"
	"github.com/Scalingo/go-scalingo/v4"
	"github.com/Scalingo/go-scalingo/v4/debug"
)

var (
	connIDGenerator = make(chan int)
	defaultPort     = 10000
	defaultBind     = "127.0.0.1"
)

type TunnelOpts struct {
	App       string
	DBEnvVar  string
	Identity  string
	Bind      string
	Port      int
	Reconnect bool
}

func Tunnel(opts TunnelOpts) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	region, err := config.GetRegion(config.C, config.C.ScalingoRegion, config.GetRegionOpts{})
	if err != nil {
		return errgo.Notef(err, "fail to retrieve region information")
	}
	sshhost := region.SSH

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
	fmt.Fprintf(os.Stderr, "Building tunnel to %s\n", dbUrl.Host)

	// Just test the connection
	client, _, err := netssh.Connect(netssh.ConnectOpts{
		Host:     sshhost,
		Identity: opts.Identity,
	})
	if err != nil {
		if err == netssh.ErrNoAuthSucceed {
			return errgo.Notef(err, "please use the flag '-i /path/to/private/key' to specify your private key")
		}
		return errgo.Notef(err, "fail to connect to SSH server")
	}
	client.Close()

	if opts.Port == 0 {
		opts.Port = defaultPort
	}
	if opts.Bind == "" {
		opts.Bind = defaultBind
	}

	var tcpAddr *net.TCPAddr
	var sock *net.TCPListener
	for {
		tcpAddr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", opts.Bind, opts.Port))
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
	fmt.Fprintln(os.Stderr, "You can access your database on:")
	fmt.Printf("%v\n", sock.Addr())

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

		go func() {
			for {
				var client *ssh.Client
				retryConnect := true
				for retryConnect {
					// Do not reuse key since the connection to the SSH agent might be broken
					client, _, err = netssh.Connect(netssh.ConnectOpts{
						Host:     sshhost,
						Identity: opts.Identity,
					})
					if err != nil {
						fmt.Println("Fail to reconnect, waiting 10 seconds...")
						time.Sleep(10 * time.Second)
					} else {
						retryConnect = false
					}
				}

				err = handleConnToTunnel(client, dbUrl, connToTunnel, errs)
				if err != nil {
					debug.Println("Error happened in tunnel", err)
					if !opts.Reconnect {
						errs <- err
						return
					}
				}
				// If err is nil, the connection has been closed normally, close the routine
				if err == nil {
					return
				}
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
		if err != stdio.EOF {
			errs <- err
		}
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

		// If the connection has been closed by the CLIENT, we must stop here and return a nil error
		// If the connection has been closed by the SERVER, we must return a errTimeout and retry the connection
		clientClosedConnectionTest := fmt.Sprintf("%s->%s: use of closed network connection", sock.LocalAddr(), sock.RemoteAddr()) // Golang error checking 101 <3
		if strings.Contains(err.Error(), clientClosedConnectionTest) {
			return nil
		}

		return err
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
