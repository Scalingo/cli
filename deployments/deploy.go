package deployments

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo"
	scalingoio "github.com/Scalingo/go-scalingo/io"

	"gopkg.in/errgo.v1"
)

type DeployRes struct {
	Deployment *scalingo.Deployment `json:"deployment"`
}

func Deploy(app, archivePath, gitRef string) error {
	c := config.ScalingoClient()

	var err error
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

	params := &scalingo.DeploymentArchiveParams{
		SourceURL: archiveURL,
	}
	// TODO gitRef cannot be anything. It is used in the docker tag image. For example, it cannot
	// start with a dash
	if strings.TrimSpace(gitRef) != "" {
		params.GitRef = &gitRef
	}
	res, err := c.DeploymentArchive(app, params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		return errgo.Newf("fail to deploy the archive: %s", res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	deployRes := &DeployRes{}
	if err = json.Unmarshal(body, &deployRes); err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	err = Stream(&StreamOpts{
		AppName:      app,
		DeploymentID: deployRes.Deployment.ID,
	})
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	return nil
}

func uploadArchivePath(c *scalingo.Client, archivePath string) (string, error) {
	scalingoio.Status("Uploading archive to Scalingo...")

	sources, err := c.SourcesCreate()
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}
	archiveBytes, err := ioutil.ReadFile(archivePath)
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}

	res, err := uploadArchiveBytes(sources.UploadURL, archiveBytes)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errgo.Newf("wrong status code after upload %s", res.Status)
	}

	return sources.DownloadURL, nil
}

func uploadArchiveBytes(uploadURL string, archiveBytes []byte) (*http.Response, error) {
	scalingoio.Status("Uploading archive to Scalingo...")
	req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(archiveBytes))
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	req.Header.Add("Content-Type", mime.TypeByExtension(".gz"))
	req.ContentLength = int64(len(archiveBytes))

	return http.DefaultClient.Do(req)
}
