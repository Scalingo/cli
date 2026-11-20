package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestProjectsAddCommand(t *testing.T) {
	tests := map[string]struct {
		args          []string
		expectedError string
		expectedHelp  bool
	}{
		"it should refuse an empty organization slug": {
			args:          []string{"--org=", "my-project"},
			expectedError: "--org expects an organization slug",
			expectedHelp:  true,
		},
		"it should refuse --default with --org": {
			args:          []string{"--org", "my-org-1", "--default", "my-project"},
			expectedError: "--default can only be used when creating a personal project",
		},
		"it should refuse --default=true with --org": {
			args:          []string{"--org", "my-org-1", "--default=true", "my-project"},
			expectedError: "--default can only be used when creating a personal project",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var output bytes.Buffer
			command := projectsAddCommand
			app := cli.Command{
				Name:      "scalingo",
				Commands:  []*cli.Command{&command},
				Writer:    &output,
				ErrWriter: &output,
			}

			err := app.Run(t.Context(), append([]string{"scalingo", "projects-add"}, test.args...))

			require.EqualError(t, err, test.expectedError)
			if test.expectedHelp {
				assert.Contains(t, output.String(), "project-name")
				assert.Contains(t, output.String(), "--org")
			}
		})
	}
}
