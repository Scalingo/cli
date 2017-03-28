package deployments

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/go-scalingo/io"

	"gopkg.in/errgo.v1"
)

type DeployRes struct {
	Deployment *scalingo.Deployment `json:"deployment"`
}

func Deploy(app, archivePath, gitRef string) error {
	_, err := url.Parse(archivePath)
	if err != nil {
		// TODO differentiate local file and URL
		// For instance if error 400 is returned, it means that archivePath is not an URL (so maybe a
		// local file?)
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Downloading the source code")
	c := config.ScalingoClient()
	params := &scalingo.DeploymentArchiveParams{
		SourceURL: archivePath,
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

	pushRes := &DeployRes{}
	if err = json.Unmarshal(body, &pushRes); err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	err = Stream(&StreamOpts{
		AppName:      app,
		DeploymentID: pushRes.Deployment.ID,
	})
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	return nil
}
