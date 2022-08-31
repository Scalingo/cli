package cmd

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/Scalingo/go-scalingo/v5/debug"
)

var (
	appFlag = cli.StringFlag{
		Name:    "app",
		Aliases: []string{"a"},
		Value:   "<name>",
		Usage:   "Name of the current app",
	}
	addonFlag = cli.StringFlag{
		Name:  "addon",
		Value: "<addon_id>",
		Usage: "ID of the current addon",
	}
)

func addonNameFromFlags(c *cli.Context) string {
	var addonName string
	if c.String("addon") != "<addon_id>" {
		addonName = c.String("addon")
	} else if os.Getenv("SCALINGO_ADDON") != "" {
		addonName = os.Getenv("SCALINGO_ADDON")
	}

	debug.Println("[ADDON] Addon name is", addonName)
	return addonName
}
