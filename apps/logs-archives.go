package apps

import (
	"fmt"
	"strconv"

	"github.com/Scalingo/cli/config"
	"github.com/fatih/color"
	errgo "gopkg.in/errgo.v1"
)

func LogsArchives(appName string, page int) error {
	if page < 0 {
		return errgo.New("Page must be greather than 0.")
	}
	if page == 0 {
		page = 1
	}

	c := config.ScalingoClient()

	logsRes, err := c.LogsArchives(appName, page)
	if err != nil {
		return errgo.Mask(err)
	}

	for _, archive := range logsRes.Archives {
		fmt.Println(color.New(color.FgYellow).Sprint("To:   ") + archive.To)
		fmt.Println(color.New(color.FgYellow).Sprint("From: ") + archive.From)
		fmt.Println(color.New(color.FgYellow).Sprint("Size: ") + strconv.FormatInt(archive.Size, 10))
		fmt.Println(color.New(color.FgYellow).Sprint("Url:  ") + archive.Url)
		fmt.Println("-------")
	}

	if len(logsRes.Archives) == 0 {
		fmt.Println("No logs archives available.")
	}

	return nil
}
