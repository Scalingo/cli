package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescriptionRendering(t *testing.T) {
	t.Run("Render only description", func(t *testing.T) {
		description := CommandDescription{
			Description: "A really cool command",
		}.Render()

		assert.Equal(t, "A really cool command", description)
	})

	t.Run("Should render with one example", func(t *testing.T) {
		description := CommandDescription{
			Description: "A really cool command",
			Examples:    []string{"scalingo create my-app"},
		}.Render()

		expectedOutput := "A really cool command\n\nExample\n  $ scalingo create my-app"

		assert.Equal(t, expectedOutput, description)
	})

	t.Run("Should render with multiple examples", func(t *testing.T) {
		description := CommandDescription{
			Description: "A really cool command",
			Examples:    []string{"scalingo create my-app", "scalingo --region osc-fr1 create my-app"},
		}.Render()

		expectedOutput := "A really cool command\n\nExamples\n  $ scalingo create my-app\n  $ scalingo --region osc-fr1 create my-app"

		assert.Equal(t, expectedOutput, description)
	})

	t.Run("Should render a single reference", func(t *testing.T) {
		description := CommandDescription{
			Description: "A really cool command",
			SeeAlso:     []string{"related-command"},
		}.Render()

		expectedOutput := "A really cool command\n\n# See also 'related-command'"

		assert.Equal(t, expectedOutput, description)
	})

	t.Run("Should render multiple references", func(t *testing.T) {
		description := CommandDescription{
			Description: "A really cool command",
			SeeAlso:     []string{"related-command", "another-related-command"},
		}.Render()

		expectedOutput := "A really cool command\n\n# See also 'related-command' 'another-related-command'"

		assert.Equal(t, expectedOutput, description)
	})

	t.Run("Should render with multiple examples and references", func(t *testing.T) {
		description := CommandDescription{
			Description: "A really cool command",
			Examples:    []string{"scalingo create yay", "scalingo create pouet"},
			SeeAlso:     []string{"related-command", "another-related-command"},
		}.Render()

		expectedOutput := "A really cool command\n\nExamples\n  $ scalingo create yay\n  $ scalingo create pouet\n\n# See also 'related-command' 'another-related-command'"

		assert.Equal(t, expectedOutput, description)
	})
}
