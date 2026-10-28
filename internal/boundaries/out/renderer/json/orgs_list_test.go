package json

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/go-scalingo/v11"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	stdout := os.Stdout
	reader, writer, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = writer
	defer func() { os.Stdout = stdout }()

	fn()

	require.NoError(t, writer.Close())
	output, err := io.ReadAll(reader)
	require.NoError(t, err)
	return string(output)
}

func TestOrgsListRenderer_Render(t *testing.T) {
	tests := map[string]struct {
		organizations []scalingo.Organization
		expected      string
	}{
		"it should only render the ID, slug and role of the organizations": {
			organizations: []scalingo.Organization{
				{ID: "org-1", Slug: "acme", Role: "admin"},
				{ID: "org-2", Slug: "globex", Role: "member"},
			},
			expected: `{"organizations":[{"id":"org-1","slug":"acme","role":"admin"},{"id":"org-2","slug":"globex","role":"member"}]}` + "\n",
		},
		"it should render an empty list": {
			organizations: []scalingo.Organization{},
			expected:      `{"organizations":[]}` + "\n",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewOrgsList()
			r.SetData(t.Context(), test.organizations)

			output := captureStdout(t, func() {
				require.NoError(t, r.Render(t.Context()))
			})

			assert.Equal(t, test.expected, output)
		})
	}
}
