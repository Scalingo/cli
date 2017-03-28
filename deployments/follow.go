package deployments

import (
	"encoding/json"
	"fmt"
	stdio "io"
	"strings"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/debug"
	"github.com/Scalingo/go-scalingo"
	"golang.org/x/net/websocket"
	"gopkg.in/errgo.v1"
)

type StreamOpts struct {
	AppName      string
	DeploymentID string
}

type deployEvent struct {
	ID   string          `json:"id"`
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type logData struct {
	Content string `json:"content"`
}

type statusData struct {
	Content string `json:"Status"`
}

func Stream(opts *StreamOpts) error {
	// TODO Sometimes, the logs of the previous deployments show up at the begining...
	c := config.ScalingoClient()
	app, err := c.AppsShow(opts.AppName)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	debug.Println("Opening socket to: " + app.Links.DeploymentsStream)

	conn, err := c.DeploymentStream(app.Links.DeploymentsStream)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	var event deployEvent
	oldStatus := ""
	for {
		err := websocket.JSON.Receive(conn, &event)
		if err != nil {
			conn.Close()
			if err == stdio.EOF {
				debug.Println("Remote server broke the connection, reconnecting")
				for err != nil {
					conn, err = c.DeploymentStream(app.Links.DeploymentsStream)
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
				var logData logData
				err := json.Unmarshal(event.Data, &logData)
				if err != nil {
					config.C.Logger.Println(err)
				} else {
					fmt.Println("[LOG] " + strings.TrimSpace(logData.Content))
				}
			case "status":
				var statusData statusData
				err := json.Unmarshal(event.Data, &statusData)
				if err != nil {
					config.C.Logger.Println(err)
				} else {
					if oldStatus == "" {
						fmt.Println("[STATUS] New status: " + statusData.Content)
					} else {
						fmt.Println("[STATUS] New status: " + oldStatus + " →  " + statusData.Content)
					}
					oldStatus = statusData.Content
					if scalingo.IsFinishedString(statusData.Content) {
						return nil
					}
				}
			case "new":
				oldStatus = ""
				var newData map[string]*scalingo.Deployment
				err := json.Unmarshal(event.Data, &newData)
				if err != nil {
					config.C.Logger.Println(err)
				} else if newData["deployment"].ID != opts.DeploymentID {
					fmt.Println("DeploymentID != opts ID")
					// TODO Pas bon ça
					break
				} else {
					fmt.Println("[NEW] New deploy: " + newData["deployment"].ID + " from " + newData["deployment"].User.Username)
				}
			}
		}
	}
}
