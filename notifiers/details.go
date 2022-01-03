package notifiers

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v4"
)

func Details(app, ID string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	baseNotifier, err := c.NotifierByID(app, ID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	notifier := baseNotifier.Specialize()

	eventTypes, err := c.EventTypesList()
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
		[]string{"ID", notifier.GetID()},
		[]string{"Type", string(notifier.GetType())},
		[]string{"Name", notifier.GetName()},
		[]string{"Enabled", strconv.FormatBool(notifier.IsActive())},
		[]string{"Send all events", strconv.FormatBool(notifier.GetSendAllEvents())},
	}
	for _, v := range data {
		t.Append(v)
	}

	// Type data
	for key, value := range notifier.TypeDataMap() {
		t.Append([]string{strings.Title(key), fmt.Sprintf("%v", value)})
	}

	//Selected events
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
