package apps

import (
	"appsdeck/cli/api"
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"syscall"
)

func Run(app string, command []string) error {
	res, err := api.Run(app, command)
	if err != nil {
		return err
	}
	runStruct := make(map[string]interface{})
	ReadJson(res.Body, &runStruct)

	if res.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("Not authorized")
	}

	if _, ok := runStruct["attach"]; !ok {
		return fmt.Errorf("Unexpected answer from server")
	}

	res, socket, err := connectToRunServer(runStruct["attach"].(string))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Unauthorized")
	}

	sttyArgs := []string{"stty", "-echo", "raw"}
	fd := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}
	_, err = syscall.ForkExec("/bin/stty", sttyArgs, &syscall.ProcAttr{Dir: "", Files: fd})
	if err != nil {
		return err
	}

	running := true
	stdinChan := readerToChan(os.Stdin)
	socketChan := readerToChan(socket)

	for running {
		select {
		case input, ok := <-stdinChan:
			if ok {
				fmt.Fprint(socket, input)
			} else {
				running = false
			}
		case data, ok := <-socketChan:
			if ok {
				fmt.Print(data)
			} else {
				os.Stdin.Close()
				running = false
			}
		}
	}
	socket.Close()

	sttyArgs = []string{"stty", "echo", "cooked"}
	fd = []uintptr{os.Stdout.Fd(), os.Stderr.Fd()}
	syscall.ForkExec("/bin/stty", sttyArgs, &syscall.ProcAttr{Dir: "", Files: fd})
	return nil
}

func connectToRunServer(rawUrl string) (*http.Response, net.Conn, error) {
	req, err := http.NewRequest("POST", rawUrl, nil)
	if err != nil {
		return nil, nil, err
	}
	api.AddAuthToken(req)

	url, err := url.Parse(rawUrl)
	if err != nil {
		return nil, nil, err
	}

	dial, err := net.Dial("tcp", url.Host)
	if err != nil {
		return nil, nil, err
	}

	conn := httputil.NewClientConn(dial, nil)

	res, err := conn.Do(req)
	if err != httputil.ErrPersistEOF && err != nil {
		return nil, nil, err
	}
	connection, _ := conn.Hijack()
	return res, connection, nil
}

func readerToChan(in io.Reader) chan string {
	c := make(chan string, 10)
	reader := bufio.NewReader(in)
	go func() {
		for {
			r, n, err := reader.ReadRune()
			if err != nil {
				close(c)
				break
			}
			if n > 0 {
				c <- string(r)
			}
		}
	}()
	return c
}
