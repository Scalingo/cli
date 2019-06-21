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
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	resources, err := c.NotifiersList(app)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Type", "Name", "Enabled", "Send all events", "Selected events"})

	eventTypes, err := c.EventTypesList()
	if err != nil {
		return errgo.Notef(err, "fail to list event types")
	}

	for _, r := range resources {
		t.Append([]string{
			r.GetID(), string(r.GetType()), r.GetName(),
			strconv.FormatBool(r.IsActive()), strconv.FormatBool(r.GetSendAllEvents()),
			eventTypesToString(eventTypes, r.GetSelectedEventIDs()),
		})
	}
	t.Render()

	return nil
}

func eventTypesToString(eventTypes []scalingo.EventType, ids []string) (res string) {
	switch len(eventTypes) {
	case 0:
		res = ""
	case 1:
		for _, t := range eventTypes {
			if t.ID == ids[0] {
				res = t.Name
				break
			}
		}
	default:
		for _, t := range eventTypes {
			if t.ID == ids[0] {
				res = fmt.Sprintf("%s, ...", t.Name)
				break
			}
		}
	}
	return
}
