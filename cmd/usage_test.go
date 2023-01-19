package cmd

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestUsageRenderingForCommands(t *testing.T) {
	t.Run("Should render the help usage without flag or arguments", func(t *testing.T) {
		command := &cli.Command{
			HelpName: "scalingo command",
		}

		usageText, err := RenderUsageWithCommand(ScalingoUsageTextTemplate, command)

		assert.NoError(t, err)
		assert.Equal(t, "scalingo command", usageText)
	})

	t.Run("Should show mandatory flags", func(t *testing.T) {
		command := &cli.Command{
			HelpName: "scalingo command",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "mandatory", Required: true},
			},
		}

		usageText, err := RenderUsageWithCommand(ScalingoUsageTextTemplate, command)

		assert.NoError(t, err)
		assert.Equal(t, "scalingo command --mandatory", usageText)
	})

}

func RenderUsageWithCommand(textTemplate string, command *cli.Command) (string, error) {
	buf := &bytes.Buffer{}

	t, err := template.New("test").Parse(textTemplate)

	if err != nil {
		return "", err
	}

	if err := t.Execute(buf, command); err != nil {
		return "", err
	}

	return buf.String(), nil
}
