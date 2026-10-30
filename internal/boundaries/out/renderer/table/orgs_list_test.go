package table

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
	t.Run("it should render the ID, slug and role of the organizations", func(t *testing.T) {
		r := NewOrgsList()
		r.SetData(t.Context(), []scalingo.Organization{
			{ID: "org-1", Slug: "acme", Role: "admin"},
			{ID: "org-2", Slug: "globex", Role: "member"},
		})

		output := captureStdout(t, func() {
			require.NoError(t, r.Render(t.Context()))
		})

		assert.Regexp(t, `│\s+ID\s+│\s+SLUG\s+│\s+ROLE\s+│`, output)
		assert.Regexp(t, `│\s+org-1\s+│\s+acme\s+│\s+admin\s+│`, output)
		assert.Regexp(t, `│\s+org-2\s+│\s+globex\s+│\s+member\s+│`, output)
		assert.NotContains(t, output, "NAME")
	})

	t.Run("it should render a message when there are no organizations", func(t *testing.T) {
		r := NewOrgsList()
		r.SetData(t.Context(), []scalingo.Organization{})

		output := captureStdout(t, func() {
			require.NoError(t, r.Render(t.Context()))
		})

		assert.Equal(t, "You are not a member of any organization.\n", output)
	})
}
