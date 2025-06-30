package logs

import (
	"bufio"
	"context"
	"fmt"
	stdio "io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/signals"
	"github.com/Scalingo/go-scalingo/v8/debug"
)

const (
	logsMaxBufferSize = 150000 // Size of the buffer when querying logs (in lines)
)

type WSEvent struct {
	Type      string    `json:"event"`
	Log       string    `json:"log"`
	Timestamp time.Time `json:"timestamp"`
}

func Dump(ctx context.Context, logsURL string, n int, filter string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	res, err := c.Logs(ctx, logsURL, n, filter)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 || res.StatusCode == 204 {
		io.Error("There is no log for this application")
		io.Info("Ensure your application is writing to the standard output")
		return nil
	}

	// Create a buffered channel with a maximum size of the number of log lines
	// requested.  On medium to good internet connection, we are fetching lines
	// faster than we can process them.  This buffer is here to get the logs as
	// fast as possible since the request will time out after 30s.
	buffSize := n
	if buffSize > logsMaxBufferSize { // Cap the size of the buffer (to prevent high memory allocation when user specify n=1_000_000)
		buffSize = logsMaxBufferSize
	}

	// This buffered channel will be used as a buffer between the network
	// connection and our logs processing pipeline.
	buff := make(chan string, buffSize)
	// This waitgroup is used to ensure that the logs processing pipeline is
	// finished before exiting the method.
	wg := &sync.WaitGroup{}

	// Start a goroutine that will read from buffered channel and send those
	// lines to the logs processing pipeline.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for bline := range buff {
			colorizeLogs(bline)
		}
	}()

	// Ensure that all lines are printed out before exiting this method.
	defer wg.Wait()

	// Here we used bufio to read from the response because we want to easily
	// split response in lines.
	// Note: This can look like a duplicate measure with our buffered channel
	// (buff) however it's not. The reason is that ReadBytes will fill the sr
	// Reader only if the internal buffer is empty. This means that the first
	// ReadBytes will fetch 4MB of data from the connection. Then it will use
	// this internal buffer until it runs out. However if our logs processing
	// pipeline is slow, it will never query the next 4MB of data. Hence the
	// buffered channel.
	sr := bufio.NewReader(res.Body)

	for {
		// Read one line from the response
		bline, err := sr.ReadBytes('\n')

		if err != nil {
			// If there was an error, we will exit, so we can close the buffered
			// channel and let the goroutine finish its work.
			close(buff)

			if err == stdio.EOF {
				// If the error is EOF, it means that we successfully read all of the
				// response body
				return nil
			}
			// Otherwise there was an error: return it
			return errgo.Notef(err, "fail to read logs")
		}
		// Send the line to the buffer
		buff <- string(bline)
	}
}

func Stream(ctx context.Context, logsRawURL string, filter string) error {
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

	header := http.Header{}
	header.Add("Origin", fmt.Sprintf("http://scalingo-cli.local/%s", config.Version))
	conn, resp, err := websocket.DefaultDialer.DialContext(ctx, logsURLString, header)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer resp.Body.Close()

	signals.CatchQuitSignals = false
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)

	go func() {
		defer close(signals)
		<-signals
		err := conn.Close()
		if err != nil {
			debug.Println("Fail to close log websocket connection", err)
		}
	}()

	for {
		err := conn.ReadJSON(&event)
		if err != nil {
			conn.Close()
			if err == stdio.EOF {
				debug.Println("Remote server broke the connection, reconnecting")
				for err != nil {
					conn, resp, err = websocket.DefaultDialer.DialContext(ctx, logsURLString, header)
					defer resp.Body.Close()
					time.Sleep(time.Second * 1)
				}
				continue
			} else if strings.Contains(err.Error(), "use of closed network connect") {
				return nil
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

		colorID := 0
		for _, letter := range []byte(container) {
			colorID += int(letter)
		}

		if container == "router" {
			colorID += 6
			content = colorizeRouterLogs(content)
		} else {
			content = errorHighlight(content)
		}
		colorID = colorID % len(containerColors)

		fmt.Printf(
			"%s [%s] %s\n",
			color.New(color.FgYellow).Sprint(date),
			containerColors[colorID](container),
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
