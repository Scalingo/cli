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

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v8"
	"github.com/Scalingo/go-utils/errors/v2"
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
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	if strings.HasPrefix(warPath, "http://") || strings.HasPrefix(warPath, "https://") {
		warReadStream, warSize, err = getURLInfo(ctx, warPath)
		if err != nil {
			return errors.Wrapf(ctx, err, "get WAR URL info")
		}
		warFileName = appName + ".war"
	} else {
		warReadStream, warSize, warFileName, err = getFileInfo(ctx, warPath)
		if err != nil {
			return errors.Wrapf(ctx, err, "get WAR file info")
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
		return errors.New(ctx, "unknown WAR size")
	}

	// Get the sources endpoints
	sources, err := c.SourcesCreate(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "create WAR deployment source")
	}

	archiveBuffer := new(bytes.Buffer)
	gzWriter := gzip.NewWriter(archiveBuffer)
	tarWriter := tar.NewWriter(gzWriter)

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return errors.Wrapf(ctx, err, "create tarball")
	}

	_, err = io.Copy(tarWriter, warReadStream)
	if err != nil {
		return errors.Wrapf(ctx, err, "copy war content")
	}

	tarWriter.Close()
	gzWriter.Close()

	res, err := uploadArchive(ctx, sources.UploadURL, archiveBuffer, int64(archiveBuffer.Len()))
	if err != nil {
		return errors.Wrapf(ctx, err, "upload the WAR archive")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.Newf(ctx, "wrong status code after upload: %s", res.Status)
	}

	return Deploy(ctx, appName, sources.DownloadURL, gitRef, opts)
}

func getURLInfo(ctx context.Context, warPath string) (io.ReadCloser, int64, error) {
	res, err := http.Get(warPath)
	if err != nil {
		return nil, 0, errors.Wrapf(ctx, err, "get WAR file")
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

func getFileInfo(ctx context.Context, warPath string) (io.ReadCloser, int64, string, error) {
	warSize := int64(0)
	warFileName := filepath.Base(warPath)
	fi, err := os.Stat(warPath)
	if err == nil {
		// If there is an error, we just skip this header
		warSize = fi.Size()
	}

	warReadStream, err := os.Open(warPath)
	if err != nil {
		return nil, 0, "", errors.Wrapf(ctx, err, "open WAR file")
	}

	return warReadStream, warSize, warFileName, nil
}
