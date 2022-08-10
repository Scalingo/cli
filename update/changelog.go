package update

import (
	"bytes"
	"encoding/json"
	"fmt"
	stdio "io"
	"net/http"
	"strings"
	"time"

	"gopkg.in/errgo.v1"
)

const (
	latestReleaseURL = "https://api.github.com/repos/scalingo/cli/releases/latest"
)

type githubChangelog struct {
	Body string `json:"body"`
}

func ShowLastChangelog() error {
	client := http.Client{
		Timeout: 4 * time.Second,
	}

	res, err := client.Get(latestReleaseURL)
	if err != nil {
		return errgo.Notef(err, "fail to request last release")
	}
	defer res.Body.Close()
	body, err := stdio.ReadAll(res.Body)
	if err != nil {
		return errgo.Notef(err, "fail to read the request of last release")
	}

	var changelogBody githubChangelog
	err = json.NewDecoder(bytes.NewBuffer((body))).Decode(&changelogBody)
	if err != nil {
		return errgo.Notef(err, "fail to decode github tag body")
	}

	fmt.Printf("%v\n\n", strings.ReplaceAll(changelogBody.Body, "\\r\\n", "\r\n"))
	return nil
}
