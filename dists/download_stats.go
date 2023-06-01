package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Scalingo/go-utils/logger"
)

type Asset struct {
	Name          string `json:"name"`
	DownloadCount int    `json:"download_count"`
}

type Repo struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Repos []Repo

func main() {
	log := logger.Default()

	res, err := http.Get("https://api.github.com/repos/Scalingo/cli/releases")
	if err != nil {
		log.WithError(err).Error("Fail to query the CLI releases from GitHub")
		return
	}
	defer res.Body.Close()

	var repos Repos
	err = json.NewDecoder(res.Body).Decode(&repos)
	if err != nil {
		log.WithError(err).Error("")
		return
	}

	for _, repo := range repos {
		fmt.Printf("%v:\n", repo.TagName)
		for _, asset := range repo.Assets {
			fmt.Printf("%v → %v downloads\n", asset.Name, asset.DownloadCount)
		}
		fmt.Println()
	}
}
