package notifiers

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func Details(app, ID string) error {
	c := config.ScalingoClient()
	baseNotifier, err := c.NotifierByID(app, ID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	notifier := baseNotifier.Specialize()
	displayDetails(notifier)
	return nil
}

func displayDetails(notifier scalingo.DetailedNotifier) {
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
		if len(notifier.GetSelectedEvents()) <= 0 {
			t.Append([]string{"Selected events", ""})
		}
		for i, e := range notifier.GetSelectedEvents() {
			if i == 0 {
				t.Append([]string{"Selected events", e.Name})
			} else {
				t.Append([]string{"", e.Name})
			}

		}
	}
	t.Render()
}
