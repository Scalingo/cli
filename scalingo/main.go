package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/stvp/rollbar"
	cli "github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/signals"
	"github.com/Scalingo/cli/update"
	errors "github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
)

var (
	completionFlag = "--" + cli.BashCompletionFlag.Names()[0]
)

func defaultAction(c *cli.Context) error {
	for i := range os.Args {
		if os.Args[i] == completionFlag {
			if len(os.Args) > 2 {
				autocomplete.FlagsAutoComplete(c, os.Args[len(os.Args)-2])
			}

			return nil
		}
	}

	err := cmd.HelpCommand.Action(c)
	if err != nil {
		return errors.Notef(c.Context, err, "help command execution")
	}

	cmd.ShowSuggestions(c)
	return nil
}

func ScalingoAppComplete(c *cli.Context) {
	// At that point, flags have not been parsed by urfave/cli
	// ie. `scalingo -a <tab><tab>`
	// So we've to handle it ourselves
	args := os.Args[1:]
	nargs := []string{}
	for _, a := range args {
		if a != completionFlag {
			nargs = append(nargs, a)
		}
	}
	if len(nargs) == 1 && (nargs[0] == "-a" || nargs[0] == "--app") {
		autocomplete.FlagAppAutoComplete(c)
		return
	}
	if len(nargs) == 1 && (nargs[0] == "-r" || nargs[0] == "--remote") {
		autocomplete.FlagRemoteAutoComplete(c)
		return
	}

	autocomplete.DisplayFlags(c.App.Flags)

	for _, command := range c.App.Commands {
		fmt.Fprintln(c.App.Writer, command.FullName())
	}
}

func setHelpTemplate() {
	cli.AppHelpTemplate = cmd.ScalingoAppHelpTemplate
	cli.CommandHelpTemplate = cmd.ScalingoCommandHelpTemplate
}

func main() {
	log := logger.Default()
	ctx := logger.ToCtx(context.Background(), log)

	app := cli.NewApp()
	app.Name = "Scalingo Client"
	app.HelpName = "scalingo"
	app.Authors = []*cli.Author{{Name: "Scalingo Team", Email: "hello@scalingo.com"}}
	app.Usage = "Manage your apps and containers"
	app.Version = config.Version
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "addon", Value: "<addon_id>", Usage: "ID of the current addon", EnvVars: []string{"SCALINGO_ADDON"}},
		&cli.StringFlag{Name: "app", Aliases: []string{"a"}, Value: "<name>", Usage: "Name of the app", EnvVars: []string{"SCALINGO_APP"}},
		&cli.StringFlag{Name: "remote", Aliases: []string{"r"}, Value: "scalingo", Usage: "Name of the remote"},
		&cli.StringFlag{Name: "region", Value: "", Usage: "Name of the region to use"},
	}
	app.EnableBashCompletion = true
	app.BashComplete = func(c *cli.Context) {
		ScalingoAppComplete(c)
	}
	app.Action = defaultAction
	setHelpTemplate()

	cmds := cmd.NewAppCommands()
	// Commands
	for _, command := range cmds.Commands() {
		oldFunc := command.BashComplete
		command.BashComplete = func(c *cli.Context) {
			n := len(os.Args) - 2
			if n > 0 && !autocomplete.FlagsAutoComplete(c, os.Args[n]) && oldFunc != nil {
				oldFunc(c)
			}
		}
		app.Commands = append(app.Commands, command)
	}

	go signals.Handle()

	bashComplete := false
	for i := range os.Args {
		if strings.Contains(os.Args[i], "generate-bash-completion") {
			bashComplete = true
		}
	}

	if !bashComplete {
		if len(os.Args) >= 2 && os.Args[1] == cmd.UpdateCommand.Name {
			err := update.Check()
			if err != nil {
				rollbar.Error(rollbar.ERR, err)
			}
			return
		}

		defer update.Check()
	} else {
		// If we are completing stuff, disable logging
		config.C.DisableInteractive = true
	}

	if err := app.RunContext(ctx, os.Args); err != nil {
		fmt.Println("Fail to run command:", err)
	}
}
