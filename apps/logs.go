package apps

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/logs"
	"github.com/Scalingo/go-scalingo/v4"
	"github.com/Scalingo/go-scalingo/v4/debug"
	"gopkg.in/errgo.v1"
)

type WSEvent struct {
	Type      string    `json:"event"`
	Log       string    `json:"log"`
	Timestamp time.Time `json:"timestamp"`
}

type LogsRes struct {
	LogsURL string        `json:"logs_url"`
	App     *scalingo.App `json:"app"`
}

func Logs(appName string, stream bool, n int, filter string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	err = checkFilter(c, appName, filter)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	res, err := c.LogsURL(appName)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errgo.Newf("fail to query logs: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	debug.Println("[API-Response] ", string(body))

	logsRes := &LogsRes{}
	if err = json.Unmarshal(body, &logsRes); err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	if err = logs.Dump(logsRes.LogsURL, n, filter); err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	if stream {
		if err = logs.Stream(logsRes.LogsURL, filter); err != nil {
			return errgo.Mask(err, errgo.Any)
		}
	}
	return nil
}

func checkFilter(c *scalingo.Client, appName string, filter string) error {
	if filter == "" {
		return nil
	}

	if strings.HasPrefix(filter, "one-off-") || strings.HasPrefix(filter, "postdeploy-") {
		return nil
	}

	processes, err := c.AppsPs(appName)
	if err != nil {
		return errgo.Mask(err)
	}

	filters := strings.Split(filter, "|")
	for _, f := range filters {
		ctName := ""

		for _, ct := range processes {
			ctName = ct.Name
			if strings.HasPrefix(f, ctName+"-") || f == ctName {
				break
			}
		}

		if !strings.HasPrefix(f, ctName+"-") && f != ctName {
			return errgo.Newf(
				"%s is not a valid container filter\n\nEXAMPLES:\n"+
					"\"scalingo logs -F web\": logs of every web containers\n"+
					"\"scalingo logs -F web-1\": logs of web container 1\n"+
					"\"scalingo logs -F web|worker\": logs of every web and worker containers\n",
				f)
		}
	}

	return nil
}
