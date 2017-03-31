package deployments

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
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

	// The tar writer will write to the pipe. At the other end of the pipe we have the sources URL to
	// upload to Scalingo.
	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()
	gzWriter := gzip.NewWriter(pipeWriter)
	tarWriter := tar.NewWriter(gzWriter)
	defer gzWriter.Close()
	defer tarWriter.Close()

	//tarErrorChannel := make(chan error)
	go func() {
		defer pipeWriter.Close()
		defer gzWriter.Close()
		defer tarWriter.Close()

		err = tarWriter.WriteHeader(header)
		if err != nil {
			//tarErrorChannel <- errgo.Mask(err, errgo.Any)
			fmt.Println(errgo.Mask(err, errgo.Any))
			return
		}

		_, err := io.Copy(tarWriter, warReadStream)
		if err != nil {
			//tarErrorChannel <- errgo.Mask(err, errgo.Any)
			fmt.Println(errgo.Mask(err, errgo.Any))
			return
		}
		/*fmt.Println("before sending")
		tarErrorChannel <- errgo.New("biniou")
		fmt.Println("after sending")
		//close(tarErrorChannel)*/
	}()
	archiveBytes, err := ioutil.ReadAll(pipeReader)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	/*fmt.Println("before receiving")
	tarError := <-tarErrorChannel
	fmt.Println("after receiving")
	if tarError != nil {
		fmt.Println("error in tar")
		return errgo.Mask(err, errgo.Any)
	}*/

	res, err := uploadArchiveBytes(sources.UploadURL, archiveBytes)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errgo.Newf("wrong status code after upload %s", res.Status)
	}

	fmt.Printf("Archive downloadable at %s\n", sources.DownloadURL)

	return Deploy(appName, sources.DownloadURL, gitRef)
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
