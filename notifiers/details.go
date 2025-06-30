package notifiers

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v8"
)

func Details(ctx context.Context, app, ID string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	baseNotifier, err := c.NotifierByID(ctx, app, ID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	notifier := baseNotifier.Specialize()

	eventTypes, err := c.EventTypesList(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to list event types")
	}

	displayDetails(notifier, eventTypes)
	return nil
}

func displayDetails(notifier scalingo.DetailedNotifier, types []scalingo.EventType) {
	t := tablewriter.NewWriter(os.Stdout)
	// Basic data
	data := [][]string{
		{"ID", notifier.GetID()},
		{"Type", string(notifier.GetType())},
		{"Name", notifier.GetName()},
		{"Enabled", strconv.FormatBool(notifier.IsActive())},
		{"Send all events", strconv.FormatBool(notifier.GetSendAllEvents())},
	}
	for _, v := range data {
		t.Append(v)
	}

	// Type data
	caser := cases.Title(language.English)
	for key, value := range notifier.TypeDataMap() {
		t.Append([]string{caser.String(key), fmt.Sprintf("%v", value)})
	}

	// Selected events
	if !notifier.GetSendAllEvents() {
		if len(notifier.GetSelectedEventIDs()) <= 0 {
			t.Append([]string{"Selected events", ""})
		}
		for i, id := range notifier.GetSelectedEventIDs() {
			var e scalingo.EventType
			for _, t := range types {
				if t.ID == id {
					e = t
					break
				}
			}

			if i == 0 {
				t.Append([]string{"Selected events", e.Name})
			} else {
				t.Append([]string{"", e.Name})
			}
		}
	}
	t.Render()
}
