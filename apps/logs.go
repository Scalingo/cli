package apps

import (
	"encoding/json"
	"fmt"
	stdio "io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Scalingo/cli/Godeps/_workspace/src/golang.org/x/net/websocket"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/debug"
	"github.com/Scalingo/cli/io"
)

type WSEvent struct {
	Type      string    `json:"event"`
	Log       string    `json:"log"`
	Timestamp time.Time `json:"timestamp"`
}

type LogsRes struct {
	LogsURL string   `json:"logs_url"`
	App     *api.App `json:"app"`
}

func Logs(appName string, stream bool, n int, filter string) error {
	err := checkFilter(appName, filter)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	res, err := api.LogsURL(appName)
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

	if err = dumpLogs(logsRes.LogsURL, n, filter); err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	if stream {
		if err = streamLogs(logsRes.LogsURL, filter); err != nil {
			return errgo.Mask(err, errgo.Any)
		}
	}
	return nil
}

func dumpLogs(logsURL string, n int, filter string) error {
	res, err := api.Logs(logsURL, n, filter)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		io.Error("There is not log for this application.")
		return nil
	}

	_, err = stdio.Copy(os.Stdout, res.Body)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	return nil
}

func streamLogs(logsRawURL string, filter string) error {
	var (
		err    error
		buffer [2048]byte
		event  WSEvent
	)

	logsURL, err := url.Parse(logsRawURL)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	if logsURL.Scheme == "https" {
		logsURL.Scheme = "wss"
	} else {
		logsURL.Scheme = "ws"
	}

	logsURLString := fmt.Sprintf("%s&stream=true", logsURL.String())
	if filter != "" {
		logsURLString = fmt.Sprintf("%s&filter=%s", logsURLString, filter)
	}

	conn, err := websocket.Dial(logsURLString, "", "http://scalingo-cli.local/"+config.Version)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	for {
		n, err := conn.Read(buffer[:])
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
		debug.Println(string(buffer[:n]))
		err = json.Unmarshal(buffer[:n], &event)
		if err != nil {
			return errgo.Notef(err, "invalid JSON %v", string(buffer[:n]))
		}
		switch event.Type {
		case "ping":
		case "log":
			fmt.Println(strings.TrimSpace(event.Log))
		}
	}

	return nil
}

func checkFilter(appName string, filter string) error {
	if filter != "" {
		processes, err := api.AppsPs(appName)
		if err != nil {
			return errgo.Mask(err)
		}

		filters := strings.Split(filter, "|")
		for _, f := range filters {

			tmpFilter := ""
			for _, ct := range processes {

				for i := ct.Amount; i > 0 && f != tmpFilter; i-- {
					if strings.Contains(f, "-") {
						tmpFilter = fmt.Sprintf("%s-%d", ct.Name, i)
					} else {
						tmpFilter = ct.Name
					}
				}
				if tmpFilter == f {
					break
				}
			}
			if tmpFilter != f {
				return errgo.Newf("%s is not a valid container type", f)
			}
		}
	}

	return nil
}
