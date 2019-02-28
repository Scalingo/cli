package db

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/briandowns/spinner"
	"github.com/cheggaaa/pb"
	errgo "gopkg.in/errgo.v1"
)

type DownloadBackupOpts struct {
	Output string
	Silent bool
}

func DownloadBackup(app, addon, backupID string, opts DownloadBackupOpts) error {
	var fileWriter io.Writer
	var logWriter io.Writer
	writeToStdout := false

	if opts.Output == "-" {
		logWriter = os.Stderr
		fileWriter = os.Stdout
		writeToStdout = true
	} else {
		logWriter = os.Stdout
	}

	if opts.Silent {
		logWriter = ioutil.Discard
	}

	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Suffix = " Preparing download"
	spinner.Writer = logWriter
	spinner.Start()

	client := config.ScalingoClient()

	backup, err := client.BackupShow(app, addon, backupID)
	if err != nil {
		return errgo.Notef(err, "fail to get backup")
	}

	filepath := ""
	if !writeToStdout {
		filepath = fmt.Sprintf("%s.tar.gz", backup.Name)
		if opts.Output != "" {
			if isDir(opts.Output) {
				filepath = fmt.Sprintf("%s/%s.tar.gz", opts.Output, backup.Name)
			} else {
				filepath = opts.Output
			}
		}
		f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return errgo.Notef(err, "fail to open file")
		}
		fileWriter = f
	}

	downloadURL, err := client.BackupDownloadURL(app, addon, backupID)
	if err != nil {
		return errgo.Notef(err, "fail to get backup download URL")
	}

	resp, err := http.Get(downloadURL)
	if err != nil {
		return errgo.Notef(err, "fail to start download")
	}
	defer resp.Body.Close()
	spinner.Stop()

	bar := pb.New(int(backup.Size)).SetUnits(pb.U_BYTES)
	bar.Output = logWriter
	bar.Start()

	reader := bar.NewProxyReader(resp.Body)
	_, err = io.Copy(fileWriter, reader)
	if writeToStdout {
		bar.Finish()
	} else {
		bar.FinishPrint(fmt.Sprintf("===> %s", filepath))
	}

	if err != nil {
		return errgo.Notef(err, "fail to download file")
	}

	return nil
}

func isDir(path string) bool {
	a, err := os.Open(path)
	if err != nil {
		return false
	}
	s, err := a.Stat()
	if err != nil {
		return false
	}

	return s.IsDir()
}
