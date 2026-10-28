package table

import (
	"io"
	"os"
	"strings"
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

func tableCells(line string) []string {
	return strings.Fields(strings.ReplaceAll(line, "│", " "))
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

		// Lines 0, 2 and 5 are the table borders.
		lines := strings.Split(strings.TrimSpace(output), "\n")
		require.Len(t, lines, 6)
		assert.Equal(t, []string{"ID", "SLUG", "ROLE"}, tableCells(lines[1]))
		assert.Equal(t, []string{"org-1", "acme", "admin"}, tableCells(lines[3]))
		assert.Equal(t, []string{"org-2", "globex", "member"}, tableCells(lines[4]))
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
