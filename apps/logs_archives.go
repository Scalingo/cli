package apps

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fatih/color"

	"github.com/Scalingo/go-utils/errors/v3"

	"github.com/Scalingo/cli/config"
)

func LogsArchives(ctx context.Context, appName string, page int) error {
	if page < 0 {
		return errors.New(ctx, "Page must be greather than 0.")
	}
	if page == 0 {
		page = 1
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	logsRes, err := c.LogsArchives(ctx, appName, page)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get apps logs archives")
	}

	if len(logsRes.Archives) == 0 {
		fmt.Println("No logs archives available.")
		return nil
	}

	for _, archive := range logsRes.Archives {
		fmt.Println(color.New(color.FgYellow).Sprint("To:   ") + archive.To)
		fmt.Println(color.New(color.FgYellow).Sprint("From: ") + archive.From)
		fmt.Println(color.New(color.FgYellow).Sprint("Size: ") + strconv.FormatInt(archive.Size, 10))
		fmt.Println(color.New(color.FgYellow).Sprint("Url:  ") + archive.URL)
		fmt.Println("-------")
	}

	return nil
}
