package db

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/cheggaaa/pb/v3"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v8"
	"github.com/Scalingo/go-scalingo/v8/debug"
	httpclient "github.com/Scalingo/go-scalingo/v8/http"
)

type DownloadBackupOpts struct {
	Output string
	Silent bool
}

func DownloadBackup(ctx context.Context, app, addon, backupID string, opts DownloadBackupOpts) error {
	// Output management (manage -s and -o - flags)
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
		logWriter = io.Discard
	}

	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client to download a backup")
	}
	if backupID == "" {
		backups, err := client.BackupList(ctx, app, addon)
		if err != nil {
			return errgo.Notef(err, "fail to get the most recent backup")
		}
		if len(backups) == 0 {
			return errgo.New("this addon has no backup")
		}
		backupID, err = getLastSuccessfulBackup(backups)
		if err != nil {
			return errgo.Notef(err, "fail to get a successful backup")
		}
		fmt.Fprintln(logWriter, "-----> Selected the most recent successful backup")
	}

	// Start a spinner when loading metadatas
	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Suffix = " Preparing download"
	spinner.Writer = logWriter
	spinner.Start()

	// Get backup metadatas
	backup, err := client.BackupShow(ctx, app, addon, backupID)
	if err != nil {
		return errgo.Notef(err, "fail to get backup")
	}

	// Generate the filename and file writer
	filepath := ""
	if !writeToStdout { // No need to generate the filename nor the file writer if we're outputing to stdout
		filepath = fmt.Sprintf("%s.tar.gz", backup.Name) // Default filename
		if opts.Output != "" {                           // If the Output flag was defined
			if isDir(opts.Output) { // If it's a directory use the default filename in this directory
				filepath = fmt.Sprintf("%s/%s.tar.gz", opts.Output, backup.Name)
			} else { // If the output is not a directory use it as the filename
				filepath = opts.Output
			}
		}
		// Open the output file
		f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return errgo.Notef(err, "fail to open file")
		}
		fileWriter = f // Set the output file as the fileWriter
	}

	// Get the pre-signed download URL
	downloadURL, err := client.BackupDownloadURL(ctx, app, addon, backupID)
	if err != nil {
		return errgo.Notef(err, "fail to get backup download URL")
	}

	debug.Println("Temporary URL to download backup is: ", downloadURL)
	// Start the download
	resp, err := http.Get(downloadURL)
	if err != nil {
		return errgo.Notef(err, "fail to start download")
	}
	defer resp.Body.Close()

	// Stop the spinner
	spinner.Stop()

	if resp.StatusCode != http.StatusOK {
		return httpclient.NewRequestFailedError(resp, &httpclient.APIRequest{
			URL:    downloadURL,
			Method: "GET",
		})
	}

	// Start the progress bar
	bar := pb.New64(int64(backup.Size)).
		Set(pb.Bytes, true).
		SetWriter(logWriter)
	bar.Start()
	reader := bar.NewProxyReader(resp.Body) // Did I tell you this library is awesome?
	_, err = io.Copy(fileWriter, reader)
	if writeToStdout { // If we were writing the file to Stdout do not print the filepath at the end
		bar.Finish()
	} else {
		// If we were writing the backup to a file, write the file path
		bar.Finish()
		fmt.Fprintf(logWriter, "===> %s\n", filepath)
	}

	if err != nil {
		return errgo.Notef(err, "fail to download file")
	}

	return nil
}

// isDir returns true if it's a valid path to a directory, false otherwise
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

func getLastSuccessfulBackup(backups []scalingo.Backup) (string, error) {
	for _, backup := range backups {
		if backup.Status == scalingo.BackupStatusDone {
			return backup.ID, nil
		}
	}
	return "", errgo.New("can't find any successful backup")
}
