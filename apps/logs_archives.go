package apps

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func LogsArchives(appName string, page int) error {
	if page < 0 {
		return errgo.New("Page must be greather than 0.")
	}
	if page == 0 {
		page = 1
	}

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	logsRes, err := c.LogsArchives(appName, page)
	if err != nil {
		return errgo.Notef(err, "fail to get apps logs archives")
	}

	if len(logsRes.Archives) == 0 {
		fmt.Println("No logs archives available.")
		return nil
	}

	for _, archive := range logsRes.Archives {
		fmt.Println(color.New(color.FgYellow).Sprint("To:   ") + archive.To)
		fmt.Println(color.New(color.FgYellow).Sprint("From: ") + archive.From)
		fmt.Println(color.New(color.FgYellow).Sprint("Size: ") + strconv.FormatInt(archive.Size, 10))
		fmt.Println(color.New(color.FgYellow).Sprint("Url:  ") + archive.Url)
		fmt.Println("-------")
	}

	return nil
}
