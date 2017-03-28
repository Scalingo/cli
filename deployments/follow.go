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

// We can stream the deployment logs of an application, or we can stream the logs of a specific
// deployments.
// The StreamOpts.DeploymentID argument is optional.
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
	var currentDeployment *scalingo.Deployment
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
				// If we stream logs of a specific deployment and this event is not about this one
				if opts.DeploymentID != "" && (currentDeployment == nil || opts.DeploymentID != currentDeployment.ID) {
					continue
				}
				var logData logData
				err := json.Unmarshal(event.Data, &logData)
				if err != nil {
					config.C.Logger.Println(err)
					continue
				}
				fmt.Println("[LOG] " + strings.TrimSpace(logData.Content))
			case "status":
				// If we stream logs of a specific deployment and this event is not about this one
				if opts.DeploymentID != "" && (currentDeployment == nil || opts.DeploymentID != currentDeployment.ID) {
					continue
				}
				var statusData statusData
				err := json.Unmarshal(event.Data, &statusData)
				if err != nil {
					config.C.Logger.Println(err)
					continue
				}
				if oldStatus == "" {
					fmt.Println("[STATUS] New status: " + statusData.Content)
				} else {
					fmt.Println("[STATUS] New status: " + oldStatus + " →  " + statusData.Content)
				}
				oldStatus = statusData.Content
				if opts.DeploymentID != "" && scalingo.IsFinishedString(statusData.Content) {
					return nil
				}
			case "new":
				oldStatus = ""
				var newData map[string]*scalingo.Deployment
				err := json.Unmarshal(event.Data, &newData)
				if err != nil {
					config.C.Logger.Println(err)
					continue
				}
				currentDeployment = newData["deployment"]
				fmt.Println("[NEW] New deploy: " + currentDeployment.ID + " from " + currentDeployment.User.Username)
			}
		}
	}
}
