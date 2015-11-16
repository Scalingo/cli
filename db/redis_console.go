package db

import (
	"bytes"
	"io"
	"net"
	"strings"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/apps"
)

func RedisConsole(app string) error {
	redisURL, _, password, err := dbURL(app, "SCALINGO_REDIS", []string{"redis://"})
	if err != nil {
		return errgo.Mask(err)
	}

	host, port, err := net.SplitHostPort(redisURL.Host)
	if err != nil {
		return errgo.Newf("%v has an invalid host", redisURL)
	}

	opts := apps.RunOpts{
		DisplayCmd:    "redis-console " + strings.Split(host, ".")[0],
		App:           app,
		Cmd:           []string{"redis-cli", "-h", host, "-p", port, "-a", password},
		StdinCopyFunc: redisStdinCopy,
	}

	err = apps.Run(opts)
	if err != nil {
		return errgo.Newf("Fail to run redis console: %v", err)
	}

	return nil
}

func redisStdinCopy(dst io.Writer, src io.Reader) (written int64, err error) {
	buf := make([]byte, 2*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			toWrite := bytes.Replace(buf[0:nr], []byte{'\n'}, []byte{'\r', '\n'}, -1)
			nr = len(toWrite)
			nw, ew := dst.Write(toWrite)
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}
