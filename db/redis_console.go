package db

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/Scalingo/cli/apps"
)

type RedisConsoleOpts struct {
	App          string
	Size         string
	VariableName string
}

func RedisConsole(opts RedisConsoleOpts) error {
	if opts.VariableName == "" {
		opts.VariableName = "SCALINGO_REDIS"
	}
	redisURL, _, password, err := dbURL(opts.App, opts.VariableName, []string{"redis", "rediss"})
	if err != nil {
		return err
	}

	if redisURL.Scheme == "rediss" {
		return fmt.Errorf("Redis console is not available when TLS connections are enforced")
	}

	host, port, err := net.SplitHostPort(redisURL.Host)
	if err != nil {
		return fmt.Errorf("%v has an invalid host", redisURL)
	}

	runOpts := apps.RunOpts{
		DisplayCmd:    "redis-console " + strings.Split(host, ".")[0],
		App:           opts.App,
		Cmd:           []string{"dbclient-fetcher", "redis", "&&", "redis-cli", "-h", host, "-p", port, "-a", password},
		Size:          opts.Size,
		StdinCopyFunc: redisStdinCopy,
	}

	err = apps.Run(runOpts)
	if err != nil {
		return fmt.Errorf("fail to run Redis console: %v", err)
	}

	return nil
}

func redisStdinCopy(dst io.Writer, src io.Reader) (int64, error) {
	var written int64
	buf := make([]byte, 2*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			toWrite := bytes.ReplaceAll(buf[0:nr], []byte{'\n'}, []byte{'\r', '\n'})
			nr = len(toWrite)
			nw, ew := dst.Write(toWrite)
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				return written, ew
			}
			if nr != nw {
				return written, io.ErrShortWrite
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			return written, er
		}
	}
	return written, nil
}
