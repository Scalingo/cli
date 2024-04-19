package notifiers

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	scalingo "github.com/Scalingo/go-scalingo/v7"
)

func List(ctx context.Context, app string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	notifiers, err := c.NotifiersList(ctx, app)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Type", "Name", "Enabled", "Send all events", "Selected events"})

	eventTypes, err := c.EventTypesList(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to list event types")
	}

	for _, notifier := range notifiers {
		selectedEvents := "All"
		if !notifier.GetSendAllEvents() {
			selectedEvents = eventTypesToString(eventTypes, notifier.GetSelectedEventIDs())
		}
		t.Append([]string{
			notifier.GetID(), string(notifier.GetType()), notifier.GetName(),
			strconv.FormatBool(notifier.IsActive()), strconv.FormatBool(notifier.GetSendAllEvents()),
			selectedEvents,
		})
	}
	t.Render()

	return nil
}

func eventTypesToString(eventTypes []scalingo.EventType, ids []string) string {
	res := ""
	if len(ids) == 0 {
		res = ""
		return res
	}
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
	return res
}
