package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	res, err := http.Get("https://api.github.com/repos/Scalingo/cli/releases")
	if err != nil {
		log.Fatalln(err)
	}

	var repos Repos
	err = json.NewDecoder(res.Body).Decode(&repos)
	if err != nil {
		log.Fatalln(err)
	}

	for _, repo := range repos {
		fmt.Printf("%v:\n", repo.TagName)
		for _, asset := range repo.Assets {
			fmt.Printf("%v → %v downloads\n", asset.Name, asset.DownloadCount)
		}
		fmt.Println()
	}
}
