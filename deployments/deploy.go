package deployments

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v4"
	"github.com/Scalingo/go-scalingo/v4/debug"
	scalingoio "github.com/Scalingo/go-scalingo/v4/io"

	"gopkg.in/errgo.v1"
)

type DeployRes struct {
	Deployment *scalingo.Deployment `json:"deployment"`
}

func Deploy(app, archivePath, gitRef string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	var archiveURL string
	// If archivePath is a remote resource
	if strings.HasPrefix(archivePath, "http://") || strings.HasPrefix(archivePath, "https://") {
		archiveURL = archivePath
	} else { // if archivePath is a file
		archiveURL, err = uploadArchivePath(c, archivePath)
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
	}

	params := &scalingo.DeploymentsCreateParams{
		SourceURL: archiveURL,
	}

	// TODO gitRef cannot be anything. It is used in the docker tag image. For example, it cannot
	// start with a dash
	if strings.TrimSpace(gitRef) != "" {
		params.GitRef = &gitRef
	}
	deployment, err := c.DeploymentsCreate(app, params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	scalingoio.Status("Your deployment has been queued and is going to start…")

	go showQueuedWarnings(c, app, deployment.ID)

	debug.Println("Streaming deployment logs of", app, ":", deployment.ID)
	err = Stream(&StreamOpts{
		AppName:      app,
		DeploymentID: deployment.ID,
	})
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	return nil
}

func uploadArchivePath(c *scalingo.Client, archivePath string) (string, error) {
	archiveFd, err := os.OpenFile(archivePath, os.O_RDONLY, 0640)
	if err != nil {
		return "", errgo.Notef(err, "fail to open archive: %v", archivePath)
	}
	defer archiveFd.Close()

	stat, err := archiveFd.Stat()
	if err != nil {
		return "", errgo.Notef(err, "fail to stat archive: %v", archivePath)
	}

	sources, err := c.SourcesCreate()
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}

	res, err := uploadArchive(sources.UploadURL, archiveFd, stat.Size())
	if err != nil {
		return "", errgo.Notef(err, "fail to upload archive: %v", archivePath)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return "", errgo.Newf("wrong status code after upload %s, body: %s", res.Status, string(body))
	}

	return sources.DownloadURL, nil
}

func uploadArchive(uploadURL string, archiveReader io.Reader, archiveSize int64) (*http.Response, error) {
	scalingoio.Status("Uploading archive…")
	req, err := http.NewRequest("PUT", uploadURL, archiveReader)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	req.Header.Set("Content-Type", "application/x-gzip")
	req.ContentLength = archiveSize

	debug.Println("Uploading archive to ", uploadURL, "with headers", req.Header)

	return http.DefaultClient.Do(req)
}

func showQueuedWarnings(c *scalingo.Client, appID, deploymentID string) {
	for {
		time.Sleep(time.Minute)
		deployment, err := c.Deployment(appID, deploymentID)
		if err != nil {
			debug.Printf("Queued deployment watcher error: %s\n", err.Error())
		}
		if deployment.Status != scalingo.StatusQueued {
			return
		}
		scalingoio.Warning("All deployment slots of application owner are currently in use, the deployment will start as soon as one is available.")
	}
}
