package apps

import (
	"bufio"
	"encoding/json"
	"fmt"
	stdio "io"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/debug"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
	"github.com/fatih/color"
	"golang.org/x/net/websocket"
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

	if res.StatusCode == 404 || res.StatusCode == 204 {
		io.Error("There is not log for this application")
		io.Info("Ensure your application is writing to the standard output")
		return nil
	}

	sr := bufio.NewReader(res.Body)
	for {
		bline, err := sr.ReadBytes('\n')
		if err != nil {
			break
		}

		colorizeLogs(string(bline))
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
				colorizeLogs(strings.TrimSpace(event.Log))
			}
		}
	}
}

func checkFilter(appName string, filter string) error {
	if filter == "" {
		return nil
	}

	if strings.HasPrefix(filter, "one-off-") || strings.HasPrefix(filter, "postdeploy-") {
		return nil
	}

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

	return nil
}

type colorFunc func(...interface{}) string

func colorizeLogs(logs string) {
	containerColors := []colorFunc{
		color.New(color.FgBlue).SprintFunc(),
		color.New(color.FgCyan).SprintFunc(),
		color.New(color.FgGreen).SprintFunc(),
		color.New(color.FgMagenta).SprintFunc(),
		color.New(color.FgHiYellow).SprintFunc(),
		color.New(color.FgHiBlue).SprintFunc(),
		color.New(color.FgHiCyan).SprintFunc(),
		color.New(color.FgHiGreen).SprintFunc(),
		color.New(color.FgHiMagenta).SprintFunc(),
	}

	lines := strings.Split(logs, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		lineSplit := strings.Split(line, " ")
		if len(lineSplit) < 5 {
			fmt.Println(line)
			continue
		}
		content := strings.Join(lineSplit[5:], " ")

		headerSplit := lineSplit[:5]
		date := strings.Join(headerSplit[:4], " ")
		containerWithSurround := headerSplit[4]
		container := containerWithSurround[1 : len(containerWithSurround)-1]

		colorId := 0
		for _, letter := range []byte(container) {
			colorId += int(letter)
		}

		if container == "router" {
			colorId += 6
			content = colorizeRouterLogs(content)
		} else {
			content = errorHighlight(content)
		}
		colorId = colorId % len(containerColors)

		fmt.Printf(
			"%s [%s] %s\n",
			color.New(color.FgYellow).Sprint(date),
			containerColors[colorId](container),
			content,
		)
	}
}

const (
	varNameState int = iota
	equalState
	inTextState
	inStringState
)

func colorizeRouterLogs(content string) string {
	var outContent string
	var stateBeginnedAt int

	state := varNameState
	// Remember where the matching state started
	stateBeginnedAt = 0
	outContent = ""
	// Will be true if we are on the last char
	isEnd := false
	for i := 0; i < len([]rune(content)); i++ {
		c := []rune(content)[i]
		if i+1 >= len([]rune(content)) {
			isEnd = true
		}

		// Some cases can return one char back if they go too far
		switch state {
		case varNameState:
			if isEnd || (!unicode.IsLetter(c) && string(c) != "_") {
				end := i
				outContent += color.New(color.FgGreen).Sprint(content[stateBeginnedAt:end])
				state = equalState
				stateBeginnedAt = end
				i--
			}
		case equalState:
			end := i + 1
			outContent += color.New(color.FgRed).Sprint(content[stateBeginnedAt:end])
			state = inTextState
			stateBeginnedAt = end
		case inTextState:
			if !isEnd && string(c) == "\"" {
				state = inStringState
			} else if isEnd || string(c) == " " {
				end := i + 1
				outContent += color.New(color.FgWhite).Sprint(content[stateBeginnedAt:end])
				state = varNameState
				stateBeginnedAt = end
			}
		case inStringState:
			isEndAfter := i+2 >= len([]rune(content))
			if !isEndAfter && string(c) == "\\" {
				// Skip next char
				i++
			} else if isEnd || string(c) == "\"" {
				state = inTextState
				if isEnd {
					i--
				}
			}
		default:
			outContent += string(c)
		}
	}

	return outContent
}

func errorHighlight(content string) string {
	reg := regexp.MustCompile("(?i)(\\berr(or)?\\b)")
	outContent := reg.ReplaceAllString(content, color.New(color.BgRed).Sprint("$1"))

	return outContent
}
