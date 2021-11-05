package cmd

import (
	"os"

	"github.com/urfave/cli"

	"github.com/Scalingo/go-scalingo/v4/debug"
)

var (
	appFlag = cli.StringFlag{
		Name:  "app, a",
		Value: "<name>",
		Usage: "Name of the current app",
	}
	addonFlag = cli.StringFlag{
		Name:  "addon",
		Value: "<addon_id>",
		Usage: "ID of the current addon",
	}
)

func addonName(c *cli.Context) string {
	var addonName string
	if c.GlobalString("addon") != "<addon_id>" {
		addonName = c.GlobalString("addon")
	} else if c.String("addon") != "<addon_id>" {
		addonName = c.String("addon")
	} else if os.Getenv("SCALINGO_ADDON") != "" {
		addonName = os.Getenv("SCALINGO_ADDON")
	}

	debug.Println("[ADDON] Addon name is", addonName)
	return addonName
}
