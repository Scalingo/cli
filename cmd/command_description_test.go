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
}
