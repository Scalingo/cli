package cmd

import (
	"fmt"
	"os"

	"github.com/Scalingo/go-scalingo/v4/debug"
	"github.com/urfave/cli"
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
	if addonName == "" {
		fmt.Println("Unable to find the addon name, please use --addon flag.")
		os.Exit(1)
	}

	debug.Println("[ADDON] Addon name is", addonName)
	return addonName
}
