package cmd

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/config"
)

var (
	VersionCommand = cli.Command{
		Name:        "version",
		Category:    "CLI Internals",
		Usage:       "Display current version",
		Description: `Display current version`,
		Action: func(c *cli.Context) {
			fmt.Println("version:", config.Version)
		},
	}
)
