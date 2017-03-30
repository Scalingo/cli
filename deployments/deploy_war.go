package deployments

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Scalingo/cli/config"

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
	// 3. send it to the new /sources endpoint
	// 4. Use the signed URL from 3. to deploy the code

	var warReadStream io.ReadCloser
	var warSize int64
	var warFileName string
	var err error
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
	// Create the tar header
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
	} else {
		return errgo.New("Unknown WAR size")
	}

	// Get the sources endpoints
	c := config.ScalingoClient()
	sources, err := c.SourcesCreate()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	fmt.Printf("Upload archive to %s\n", sources.UploadURL)

	// The tar writer will write to the pipe. At the other end of the pipe we have the sources URL to
	// upload to Scalingo.
	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()
	gzWriter := gzip.NewWriter(pipeWriter)
	tarWriter := tar.NewWriter(gzWriter)
	defer gzWriter.Close()
	defer tarWriter.Close()

	go func() {
		defer pipeWriter.Close()
		defer gzWriter.Close()
		defer tarWriter.Close()

		err = tarWriter.WriteHeader(header)
		if err != nil {
			// TODO Use channel to get error back in the main thread
			fmt.Println("error")
			fmt.Println(errgo.Mask(err, errgo.Any))
		}

		_, err := io.Copy(tarWriter, warReadStream)
		if err != nil {
			fmt.Println("error")
			fmt.Println(errgo.Mask(err, errgo.Any))
		}
	}()
	var bodyDuplicate bytes.Buffer
	teeReader := io.TeeReader(pipeReader, &bodyDuplicate)

	req, err := http.NewRequest("PUT", sources.UploadURL, teeReader)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	req.GetBody = func() (io.ReadCloser, error) {
		r := bytes.NewReader(bodyDuplicate.Bytes())
		return ioutil.NopCloser(r), nil
	}

	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	req.Header.Add("Content-Type", mime.TypeByExtension(".gz"))
	// TODO mandatory header but impossible with full stream
	req.ContentLength = 1575

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Printf("body: %+v\n", string(body))
		return errgo.Newf("wrong status code %s", res.Status)
	}

	fmt.Printf("Archive downloadable at %s\n", sources.DownloadURL)

	//return Deploy(appName, sources.DownloadURL, gitRef)
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

func _devTGzipArchive(tarWriter *tar.Writer, warReadStream io.ReadCloser, header *tar.Header) error {
	err := tarWriter.WriteHeader(header)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	_, err = io.Copy(tarWriter, warReadStream)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	return nil
}
