package apps

import (
	"github.com/Appsdeck/appsdeck/api"
	"bufio"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func Logs(appName string, stream bool, n int) error {
	res, err := api.LogsURL(appName)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("fail to query logs: %s", res.Status)
	}

	app := App{}
	if err = json.NewDecoder(res.Body).Decode(&app); err != nil {
		return err
	}

	res, err = api.Logs(app.LogsURL, stream, n)
	if err != nil {
		return err
	}

	if !stream {
		buffer, _ := ioutil.ReadAll(res.Body)
		fmt.Println(html.UnescapeString(string(buffer)))
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
				html.UnescapeString(
					strings.TrimSpace(parsedLine[1]),
				),
			)
		}
	}
	if err != io.EOF {
		return err
	}
	return nil
}
