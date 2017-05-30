package apps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/Scalingo/cli/config"
	errgo "gopkg.in/errgo.v1"
)

type LogsItem struct {
	Url  string `json:"url"`
	From string `json:"from"`
	To   string `json:"to"`
	Size int64  `json:"size"`
}

type LogsResponse struct {
	NextCursor string     `json:"next_cursor"`
	HasMore    bool       `json:"has_more"`
	Archives   []LogsItem `json:"archives"`
}

func LogsArchives(appName string, cursor string) error {
	c := config.ScalingoClient()

	res, err := c.LogsArchives(appName, cursor)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	var logsRes = LogsResponse{}
	err = json.Unmarshal(body, &logsRes)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
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
