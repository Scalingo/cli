package apps

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/debug"
)

type LogsRes struct {
	LogsURL string   `json:"logs_url"`
	App     *api.App `json:"app"`
}

func Logs(appName string, stream bool, n int) error {
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

	res, err = api.Logs(logsRes.LogsURL, stream, n)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	if !stream {
		buffer, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(buffer))
	} else {
		return streamLogs(res)
	}
	return nil
}

func streamLogs(res *http.Response) error {
	var err error
	reader := bufio.NewReader(res.Body)
	for line, _, err := reader.ReadLine(); err == nil; line, _, err = reader.ReadLine() {
		if len(line) != 0 {
			parsedLine := strings.SplitN(string(line), ":", 2)
			if len(parsedLine) != 2 {
				// Invalid content from server, SSE should be
				// msgname: content
				// Anything else is wrong
				continue
			}
			fmt.Println(
				strings.TrimSpace(parsedLine[1]),
			)
		}
	}
	if err != io.EOF {
		return errgo.Mask(err, errgo.Any)
	}
	return nil
}
