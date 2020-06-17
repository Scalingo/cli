package log_drains

import (
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func ListAddon(app string, addonID string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	logDrainsAddonList, err := c.LogDrainsAddonList(app, addonID)
	if err != nil {
		return errgo.Notef(err, "fail to list the log drains")
	}

	t := tablewriter.NewWriter(os.Stdout)

	t.SetHeader([]string{"Addon name", "Addon plan", "URL"})
	t.SetAutoMergeCells(true)

	for _, logDrain := range logDrainsAddonList.Drains {
		t.Append([]string{
			logDrainsAddonList.Addon.Name,
			logDrainsAddonList.Addon.Plan,
			logDrain.URL,
		})
	}

	t.Render()
	return nil
}
