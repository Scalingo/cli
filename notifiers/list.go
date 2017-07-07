package notifiers

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Scalingo/cli/config"
	scalingo "github.com/Scalingo/go-scalingo"
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
	t.SetHeader([]string{"ID", "Type", "Name", "Enabled", "Send all events", "Selected events"})

	for _, r := range resources {
		t.Append([]string{
			r.GetID(), string(r.GetType()), r.GetName(),
			strconv.FormatBool(r.IsActive()), strconv.FormatBool(r.GetSendAllEvents()),
			eventTypesToString(r.GetSelectedEvents()), // r.TypeDataString()[:typeDataStringLength],
		})
	}
	t.Render()

	return nil
}

func eventTypesToString(eventTypes []scalingo.EventTypeStruct) (res string) {
	switch len(eventTypes) {
	case 0:
		res = ""
	case 1:
		res = eventTypes[0].Name
	default:
		res = fmt.Sprintf("%s, ...", eventTypes[0].Name)
	}
	return
}
