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

	"github.com/Scalingo/go-scalingo"
	"golang.org/x/net/websocket"
	"gopkg.in/errgo.v1"
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
	LogsURL string        `json:"logs_url"`
	App     *scalingo.App `json:"app"`
}

func Logs(appName string, stream bool, n int, filter string) error {
	err := checkFilter(appName, filter)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	c := config.ScalingoClient()
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
	c := config.ScalingoClient()
	res, err := c.Logs(logsURL, n, filter)
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
		err   error
		event WSEvent
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
		err := websocket.JSON.Receive(conn, &event)
		if err != nil {
			conn.Close()
			if err == stdio.EOF {
				debug.Println("Remote server broke the connection, reconnecting")
				for err != nil {
					conn, err = websocket.Dial(logsURLString, "", "http://scalingo-cli.local/"+config.Version)
					time.Sleep(time.Second * 1)
				}
				continue
			} else {
				return errgo.Mask(err, errgo.Any)
			}
		} else {
			switch event.Type {
			case "ping":
			case "log":
				fmt.Println(strings.TrimSpace(event.Log))
			}
		}
	}
}

func checkFilter(appName string, filter string) error {
	if filter != "" {
		c := config.ScalingoClient()
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
	}

	return nil
}
