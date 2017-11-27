package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Scalingo/cli/cmd"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/signals"
	"github.com/Scalingo/cli/update"
	"github.com/stvp/rollbar"
	"github.com/urfave/cli"
)

var (
	completionFlag = "--" + cli.BashCompletionFlag.GetName()
)

func DefaultAction(c *cli.Context) {
	completeMode := false

	for i := range os.Args {
		if os.Args[i] == completionFlag {
			completeMode = true
			break
		}
	}

	if !completeMode {
		cmd.HelpCommand.Action.(func(*cli.Context))(c)
		cmd.ShowSuggestions(c)
	} else {
		i := len(os.Args) - 2
		if i > 0 {
			autocomplete.FlagsAutoComplete(c, os.Args[i])
		}
	}
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
	cli.AppHelpTemplate = `NAME:
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION:
   {{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if len .Authors}}

AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
   {{range $index, $author := .Authors}}{{if $index}}
   {{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}
		 {{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

COPYRIGHT:
   {{.Copyright}}{{end}}
`
}

func main() {
	app := cli.NewApp()
	app.Name = "Scalingo Client"
	app.Author = "Scalingo Team"
	app.Email = "hello@scalingo.com"
	app.Usage = "Manage your apps and containers"
	app.Version = config.Version
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "app, a", Value: "<name>", Usage: "Name of the app", EnvVar: "SCALINGO_APP"},
		cli.StringFlag{Name: "remote, r", Value: "scalingo", Usage: "Name of the remote", EnvVar: ""},
	}
	app.EnableBashCompletion = true
	app.BashComplete = func(c *cli.Context) {
		ScalingoAppComplete(c)
	}
	app.Action = DefaultAction
	setHelpTemplate()

	// Commands
	for _, command := range cmd.Commands {
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
		} else {
			defer update.Check()
		}
	} else {
		// If we are completing stuff, disable logging
		config.C.DisableInteractive = true
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("Fail to run command:", err)
	}
}
