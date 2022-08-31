package events

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v5"
)

type DisplayTimelineOpts struct {
	DisplayAppName bool
}

func DisplayTimeline(events scalingo.Events, pagination scalingo.PaginationMeta, opts DisplayTimelineOpts) error {
	longestEventName := 0
	longestAppName := 0
	for _, event := range events {
		if len(event.GetEvent().Type) > longestEventName {
			longestEventName = len(event.GetEvent().Type)
		}
		if len(event.GetEvent().AppName) > longestAppName {
			longestAppName = len(event.GetEvent().AppName)
		}
	}

	for _, event := range events {
		t := event.PrintableType()
		if len(t) < longestEventName {
			for len(t) != longestEventName {
				t += " "
			}
		}

		app := event.GetEvent().AppName
		if opts.DisplayAppName && len(app) > 0 {
			if len(app) < longestAppName {
				for len(app) != longestAppName {
					app += " "
				}
			}

			fmt.Printf(
				"* %s - %s - %s - %s %s\n",
				io.Yellow(event.When()),
				io.Green(t),
				io.LightGray(app),
				event.String(),
				io.BoldBlue(
					fmt.Sprintf("<%s>", event.Who()),
				),
			)
		} else {
			fmt.Printf(
				"* %s - %s - %s %s\n",
				io.Yellow(event.When()),
				io.Green(t),
				event.String(),
				io.BoldBlue(
					fmt.Sprintf("<%s>", event.Who()),
				),
			)
		}
	}
	fmt.Fprintln(os.Stderr, io.Gray(fmt.Sprintf("Page: %d, Last Page: %d", pagination.CurrentPage, pagination.TotalPages)))
	return nil
}
