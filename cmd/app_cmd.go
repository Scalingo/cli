package cmd

import (
	"github.com/Scalingo/codegangsta-cli"
)

var (
	appFlag = cli.StringFlag{Name: "app, a", Value: "<name>", Usage: "Name of the current app"}
)
