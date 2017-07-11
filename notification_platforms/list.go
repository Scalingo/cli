package notification_platforms

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/config"
	scalingo "github.com/Scalingo/go-scalingo"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func List() error {
	c := config.ScalingoClient()
	resources, err := c.NotificationPlatformsList()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Name"})

	for _, r := range resources {
		t.Append([]string{r.Name})
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
