package cmd

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v7/debug"
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

// exitIfMissing is optional. Set to true to show a message requesting for the --addon flag.
func addonUUIDFromFlags(c *cli.Context, app string, exitIfMissing ...bool) string {
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

	if addonName == "" && len(exitIfMissing) > 0 && exitIfMissing[0] {
		io.Error("Unable to find the addon name, please use --addon flag.")
		os.Exit(1)
	}
	if addonName == "" {
		return ""
	}

	var addonUUID string
	var err error
	addonUUID, err = utils.GetAddonUUIDFromType(c.Context, app, addonName)
	if err != nil {
		io.Error("Unable to get the addon UUID based on its type:", err)
		os.Exit(1)
	}
	debug.Println("[ADDON] Addon name is", addonName)
	return addonUUID
}

func regionNameFromFlags(c *cli.Context) string {
	for _, cliContext := range c.Lineage() {
		if cliContext.String("region") != "" {
			return cliContext.String("region")
		}
	}
	return ""
}
