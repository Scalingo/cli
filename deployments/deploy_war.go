package deployments

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v6"
)

type DeployWarRes struct {
	Deployment *scalingo.Deployment `json:"deployment"`
}

func DeployWar(ctx context.Context, appName, warPath, gitRef string, opts DeployOpts) error {
	var warReadStream io.ReadCloser

	var warSize int64
	var warFileName string

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	if strings.HasPrefix(warPath, "http://") || strings.HasPrefix(warPath, "https://") {
		warReadStream, warSize, err = getURLInfo(warPath)
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
		warFileName = appName + ".war"
	} else {
		warReadStream, warSize, warFileName, err = getFileInfo(warPath)
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
	}
	defer warReadStream.Close()
	// Create the tar header
	header := &tar.Header{
		Name:       fmt.Sprintf("%s/%s", appName, warFileName),
		Typeflag:   tar.TypeReg, // Is a regular file
		Mode:       0640,
		ModTime:    time.Now(),
		AccessTime: time.Now(),
		ChangeTime: time.Now(),
	}
	if warSize != 0 {
		header.Size = warSize
	} else {
		return errgo.New("Unknown WAR size")
	}

	// Get the sources endpoints
	sources, err := c.SourcesCreate(ctx)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	archiveBuffer := new(bytes.Buffer)
	gzWriter := gzip.NewWriter(archiveBuffer)
	tarWriter := tar.NewWriter(gzWriter)

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return errgo.Notef(err, "fail to create tarball")
	}

	_, err = io.Copy(tarWriter, warReadStream)
	if err != nil {
		return errgo.Notef(err, "fail to copy war content")
	}

	tarWriter.Close()
	gzWriter.Close()

	res, err := uploadArchive(sources.UploadURL, archiveBuffer, int64(archiveBuffer.Len()))
	if err != nil {
		return errgo.Notef(err, "fail to upload the WAR archive")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errgo.Newf("wrong status code after upload %s", res.Status)
	}

	return Deploy(ctx, appName, sources.DownloadURL, gitRef, opts)
}

func getURLInfo(warPath string) (io.ReadCloser, int64, error) {
	res, err := http.Get(warPath)
	if err != nil {
		return nil, 0, errgo.Mask(err, errgo.Any)
	}

	warSize := int64(0)
	warReadStream := res.Body
	if res.Header.Get("Content-Length") != "" {
		i, err := strconv.ParseInt(res.Header.Get("Content-Length"), 10, 64)
		if err == nil {
			// If there is an error, we just skip this header
			warSize = i
		}
	}

	return warReadStream, warSize, nil
}

func getFileInfo(warPath string) (io.ReadCloser, int64, string, error) {
	warSize := int64(0)
	warFileName := filepath.Base(warPath)
	fi, err := os.Stat(warPath)
	if err == nil {
		// If there is an error, we just skip this header
		warSize = fi.Size()
	}

	warReadStream, err := os.Open(warPath)
	if err != nil {
		return nil, 0, "", errgo.Mask(err, errgo.Any)
	}

	return warReadStream, warSize, warFileName, nil
}
