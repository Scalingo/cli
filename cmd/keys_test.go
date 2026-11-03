package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestAddSSHKeyCommandRejectsInvalidArgumentCounts(t *testing.T) {
	tests := map[string][]string{
		"missing both arguments": nil,
		"missing key path":       {"key-name"},
		"extra argument":         {"key-name", "path-to-key", "extra"},
	}

	for name, args := range tests {
		t.Run(name, func(t *testing.T) {
			var output bytes.Buffer
			command := addSSHKeyCommand
			app := cli.Command{
				Name:      "scalingo",
				Commands:  []*cli.Command{&command},
				Writer:    &output,
				ErrWriter: &output,
			}

			err := app.Run(context.Background(), append([]string{"scalingo", "keys-add"}, args...))

			require.EqualError(t, err, "keys-add expects exactly 2 arguments: key-name path-to-key")
			assert.Contains(t, output.String(), "USAGE:")
			assert.Contains(t, output.String(), "key-name path-to-key")
		})
	}
}
