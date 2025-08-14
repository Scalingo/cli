package main

import (
	"context"
	"fmt"
	"net/mail"
	"os"
	"strings"

	"github.com/stvp/rollbar"
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/signals"
	"github.com/Scalingo/cli/update"
	"github.com/Scalingo/go-scalingo/v8/debug"
	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
)

var (
	completionFlag = "--" + cli.GenerateShellCompletionFlag.Names()[0]
)

func defaultAction(ctx context.Context, c *cli.Command) error {
	for i := range os.Args {
		if os.Args[i] == completionFlag {
			if len(os.Args) > 2 {
				autocomplete.FlagsAutoComplete(ctx, os.Args[len(os.Args)-2])
			}

			return nil
		}
	}

	err := cmd.HelpCommand.Action(ctx, c)
	if err != nil {
		return errors.Wrapf(ctx, err, "help command execution")
	}

	cmd.ShowSuggestions(c)

	return nil
}

func ScalingoAppComplete(ctx context.Context, c *cli.Command) {
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
		autocomplete.FlagAppAutoComplete(ctx)
		return
	}
	if len(nargs) == 1 && (nargs[0] == "-r" || nargs[0] == "--remote") {
		autocomplete.FlagRemoteAutoComplete()
		return
	}

	autocomplete.DisplayFlags(c.Flags)

	for _, command := range c.Commands {
		fmt.Fprintln(c.Writer, command.FullName())
	}
}

func setHelpTemplate() {
	cli.RootCommandHelpTemplate = cmd.ScalingoAppHelpTemplate
	cli.CommandHelpTemplate = cmd.ScalingoCommandHelpTemplate
}

func main() {
	log := logger.Default()
	ctx := logger.ToCtx(context.Background(), log)

	app := cli.Command{}
	app.Name = "Scalingo Client"
	app.Authors = []any{mail.Address{Name: "Scalingo Team", Address: "hello@scalingo.com"}}
	app.Usage = "Manage your apps and containers"
	app.Version = config.Version
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "addon", Value: "<addon_id>", Usage: "ID of the current addon", Sources: cli.EnvVars("SCALINGO_ADDON")},
		&cli.StringFlag{Name: "app", Aliases: []string{"a"}, Value: "<name>", Usage: "Name of the app", Sources: cli.EnvVars("SCALINGO_APP")},
		&cli.StringFlag{Name: "remote", Aliases: []string{"r"}, Value: "scalingo", Usage: "Name of the remote"},
		&cli.StringFlag{Name: "region", Value: "", Usage: "Name of the region to use"},
	}
	app.EnableShellCompletion = true
	app.ShellComplete = func(ctx context.Context, c *cli.Command) {
		ScalingoAppComplete(ctx, c)
	}
	app.Action = defaultAction
	setHelpTemplate()

	cmds := cmd.NewAppCommands()
	// Commands
	for _, command := range cmds.Commands() {
		oldFunc := command.ShellComplete
		command.ShellComplete = func(ctx context.Context, c *cli.Command) {
			n := len(os.Args) - 2
			if n > 0 && !autocomplete.FlagsAutoComplete(ctx, os.Args[n]) && oldFunc != nil {
				oldFunc(ctx, c)
			}
		}
		app.Commands = append(app.Commands, command)
	}

	go signals.Handle()

	bashComplete := false
	for i := range os.Args {
		if strings.Contains(os.Args[i], "generate-shell-completion") {
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
	} else {
		// If we are completing stuff, disable logging
		config.C.DisableInteractive = true
	}

	err := app.Run(ctx, os.Args)
	if err != nil {
		fmt.Println("Fail to run command:", err)
	}

	// Do not show update check during autocomplete
	if bashComplete {
		debug.Println("Do not check available update during autocompletion")
	} else {
		// We want to display to the user if a new version is available
		// Whatever the success of the execution of their command is.
		updateCheckErr := update.Check()
		if updateCheckErr != nil {
			debug.Printf("Failed to check if executable should be updated: %v\n", updateCheckErr)
		}
	}

	if err != nil {
		os.Exit(1)
	}
}
