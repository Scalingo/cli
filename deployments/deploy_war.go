package deployments

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Scalingo/go-scalingo"

	errgo "gopkg.in/errgo.v1"
)

type DeployWarRes struct {
	Deployment *scalingo.Deployment `json:"deployment"`
}

func DeployWar(appName, warPath, gitRef string) error {
	// TODO algo
	// 1. If file
	// 1.b. if url, download the archive and go to 2
	// 2. insert into a tgz (inside a folder)
	// 3. send it to the new /source endpoint
	// 4. Use the signed URL from 3. to deploy the code

	// TODO out should be the /sources endpoint
	out, err := os.Create(appName + ".tgz")
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer out.Close()

	var warReadStream io.ReadCloser
	var warSize int64
	var warFileName string
	if strings.HasPrefix(warPath, "http://") || strings.HasPrefix(warPath, "https://") {
		fmt.Println("If is URL, download the WAR")
		warReadStream, warSize, err = getURLInfo(warPath)
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
		warFileName = appName
	} else {
		fmt.Println("If is file, just archive in tgz")
		warReadStream, warSize, warFileName, err = getFileInfo(warPath)
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
	}
	defer warReadStream.Close()

	gzWriter := gzip.NewWriter(out)
	defer gzWriter.Close()
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	header := &tar.Header{
		Name:       fmt.Sprintf("%s/%s.war", appName, warFileName),
		Typeflag:   tar.TypeReg, // Is a regular file
		Mode:       0644,
		ModTime:    time.Now(),
		AccessTime: time.Now(),
		ChangeTime: time.Now(),
	}
	if warSize != 0 {
		// TODO It is mandatory. What to do if we cannot find the size
		header.Size = warSize
	}
	_ = warSize
	err = tarWriter.WriteHeader(header)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	n, err := io.Copy(tarWriter, warReadStream)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	fmt.Printf("copied bytes: %d\n", n)

	/*b, err := ioutil.ReadAll(warReadStream)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	_, err = tarWriter.Write(b)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}*/

	/*c := config.ScalingoClient()
	_ = c

	err = Stream(&StreamOpts{
		AppName: appName,
		//DeploymentID: deployRes.Deployment.ID,
	})
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}*/

	return nil
}

func getURLInfo(warPath string) (warReadStream io.ReadCloser, warSize int64, err error) {
	res, err := http.Get(warPath)
	if err != nil {
		err = errgo.Mask(err, errgo.Any)
		return
	}
	warReadStream = res.Body
	if res.Header.Get("Content-Length") != "" {
		i, err := strconv.ParseInt(res.Header.Get("Content-Length"), 10, 64)
		if err == nil {
			// If there is an error, we just skip this header
			warSize = i
		}
	}
	return
}

func getFileInfo(warPath string) (warReadStream io.ReadCloser, warSize int64, warFileName string, err error) {
	warFileName = filepath.Base(warPath)
	fi, err := os.Stat(warPath)
	if err == nil {
		// If there is an error, we just skip this header
		warSize = fi.Size()
	}
	warReadStream, err = os.Open(warPath)
	if err != nil {
		err = errgo.Mask(err, errgo.Any)
		return
	}
	return
}
