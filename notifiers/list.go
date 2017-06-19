package notifiers

import (
	"os"
	"strconv"
	"strings"

	"github.com/Scalingo/cli/config"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func List(app string) error {
	c := config.ScalingoClient()
	resources, err := c.NotifiersList(app)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Type", "Name", "Enabled", "Send all events", "Selected events", "Type data"})

	for _, r := range resources {
		// TODO: remove this
		typeDataStringLength := 10
		if len(r.TypeDataString()) < 10 {
			typeDataStringLength = len(r.TypeDataString())
		}

		t.Append([]string{
			r.GetID(), string(r.GetType()), r.GetName(),
			strconv.FormatBool(r.IsActive()), strconv.FormatBool(r.GetSendAllEvents()),
			strings.Join(r.GetSelectedEvents(), "\n"), r.TypeDataString()[:typeDataStringLength],
		})
	}
	t.Render()

	return nil
}
