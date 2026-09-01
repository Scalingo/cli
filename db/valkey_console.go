package db

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/go-utils/errors/v3"
)

type ValkeyConsoleOpts struct {
	App          string
	Size         string
	VariableName string
}

func ValkeyConsole(ctx context.Context, opts ValkeyConsoleOpts) error {
	if opts.VariableName == "" {
		opts.VariableName = "SCALINGO_VALKEY"
	}
	valkeyURL, _, password, err := dbURL(ctx, opts.App, opts.VariableName, []string{"valkey", "valkeys"})
	if err != nil {
		return errors.Wrapf(ctx, err, "resolve Valkey URL from %s", opts.VariableName)
	}

	if valkeyURL.Scheme == "valkeys" {
		return errors.New(ctx, "Valkey console is not available when TLS connections are enforced")
	}

	host, port, err := net.SplitHostPort(valkeyURL.Host)
	if err != nil {
		return fmt.Errorf("%v has an invalid host", valkeyURL)
	}

	err = apps.Run(ctx, apps.RunOpts{
		DisplayCmd:    "valkey-console " + strings.Split(host, ".")[0],
		App:           opts.App,
		Cmd:           []string{dbClientFetcher, "valkey", "&&", "valkey-cli", "-h", host, "-p", port, "-a", password},
		Size:          opts.Size,
		StdinCopyFunc: valkeyStdinCopy,
	})
	if err != nil {
		return fmt.Errorf("run Valkey console: %v", err)
	}

	return nil
}

func valkeyStdinCopy(dst io.Writer, src io.Reader) (int64, error) {
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
