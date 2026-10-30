package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/internal/boundaries/out/renderer"
	rendererjson "github.com/Scalingo/cli/internal/boundaries/out/renderer/json"
	renderertable "github.com/Scalingo/cli/internal/boundaries/out/renderer/table"
	"github.com/Scalingo/cli/orgs"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

var (
	orgsCommand = cli.Command{
		Name:     "orgs",
		Category: "Global",
		Usage:    "List your organizations",
		Description: CommandDescription{
			Description: "List the organizations you are a member of, with your role in each of them",
			Examples:    []string{"scalingo orgs", "scalingo orgs --format json"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			format := renderer.Format(c.String("format"))

			var orgsRenderer renderer.Renderer[[]scalingo.Organization]
			switch format {
			case renderer.FormatTable:
				orgsRenderer = renderertable.NewOrgsList()
			case renderer.FormatJSON:
				orgsRenderer = rendererjson.NewOrgsList()
			default:
				errorQuitWithHelpMessage(ctx, errors.Newf(ctx, "invalid format '%v'", format), c, "orgs")
			}

			err := orgs.List(ctx, orgsRenderer)
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "orgs")
		},
	}
)
