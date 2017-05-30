package apps

import (
	"fmt"
	"strconv"

	"github.com/Scalingo/cli/config"
	errgo "gopkg.in/errgo.v1"
)

func LogsArchives(appName string, cursor string) error {
	c := config.ScalingoClient()

	logsRes, err := c.LogsArchives(appName, cursor)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Println("-------")

	for _, archive := range logsRes.Archives {
		fmt.Println("| To:   " + archive.To)
		fmt.Println("| From: " + archive.From)
		fmt.Println("| Size: " + strconv.FormatInt(archive.Size, 10))
		fmt.Println("| Url:  " + archive.Url)
		fmt.Println("-------")
	}

	if !logsRes.HasMore && len(logsRes.Archives) == 0 {
		fmt.Println("No logs archives available for this app.")
	} else if logsRes.HasMore {
		fmt.Println("Next page cursor: " + logsRes.NextCursor)
	} else {
		fmt.Println("Nothing more available.")
	}

	return nil
}
