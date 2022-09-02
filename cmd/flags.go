package cmd

import (
	"os"

	"github.com/Scalingo/go-scalingo/v4/debug"
	"github.com/urfave/cli/v2"
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

	for _, cliContext := range c.Lineage() {
		if cliContext.String("addon") != "<addon_id>" {
			addonName = cliContext.String("addon")
			break
		}
	}

	if addonName == "" && os.Getenv("SCALINGO_ADDON") != "" {
		addonName = os.Getenv("SCALINGO_ADDON")
	}

	debug.Println("[ADDON] Addon name is", addonName)
	return addonName
}
