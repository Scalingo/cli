package deployments

import (
	"context"
	"io"
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

type DeployOpts struct {
	NoFollow bool
}

func Deploy(ctx context.Context, app, archivePath, gitRef string, opts DeployOpts) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	var archiveURL string
	// If archivePath is a remote resource
	if strings.HasPrefix(archivePath, "http://") || strings.HasPrefix(archivePath, "https://") {
		archiveURL = archivePath
	} else { // if archivePath is a file
		archiveURL, err = uploadArchivePath(ctx, client, archivePath)
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
	}

	params := &scalingo.DeploymentsCreateParams{
		SourceURL: archiveURL,
	}

	// gitRef cannot be anything. It is used in the docker tag image. For example, it cannot
	// start with a dash
	if strings.TrimSpace(gitRef) != "" {
		params.GitRef = &gitRef
	}
	deployment, err := client.DeploymentsCreate(ctx, app, params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	scalingoio.Status("Your deployment has been queued and is going to start…")

	if opts.NoFollow {
		scalingoio.Statusf("The no-follow flag is passed. You can check deployment logs with scalingo --app %s deployment-follow", app)
		return nil
	}

	go showQueuedWarnings(ctx, client, app, deployment.ID)
	debug.Println("Streaming deployment logs of", app, ":", deployment.ID)
	err = Stream(ctx, &StreamOpts{
		AppName:      app,
		DeploymentID: deployment.ID,
	})
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	return nil
}

func uploadArchivePath(ctx context.Context, client *scalingo.Client, archivePath string) (string, error) {
	archiveFd, err := os.OpenFile(archivePath, os.O_RDONLY, 0640)
	if err != nil {
		return "", errgo.Notef(err, "fail to open archive: %v", archivePath)
	}
	defer archiveFd.Close()

	stat, err := archiveFd.Stat()
	if err != nil {
		return "", errgo.Notef(err, "fail to stat archive: %v", archivePath)
	}

	sources, err := client.SourcesCreate(ctx)
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}

	res, err := uploadArchive(sources.UploadURL, archiveFd, stat.Size())
	if err != nil {
		return "", errgo.Notef(err, "fail to upload archive: %v", archivePath)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
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

func showQueuedWarnings(ctx context.Context, client *scalingo.Client, appID, deploymentID string) {
	for {
		time.Sleep(time.Minute)
		deployment, err := client.Deployment(ctx, appID, deploymentID)
		if err != nil {
			debug.Printf("Queued deployment watcher error: %s\n", err.Error())
		}
		if deployment.Status != scalingo.StatusQueued {
			return
		}
		scalingoio.Warning("All deployment slots of application owner are currently in use, the deployment will start as soon as one is available.")
	}
}
