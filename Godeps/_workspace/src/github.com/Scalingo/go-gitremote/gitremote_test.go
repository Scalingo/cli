package gitremote

import (
	"reflect"
	"testing"
)

func TestListFromConfigContent1Remote(t *testing.T) {
	remotes, err := ListFromConfigContent(config1Remote)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(remotes) != 1 {
		t.Fatalf("Expeced 1 remote, got %v", len(remotes))
	}

	remote := &Remote{
		Name:  "origin",
		URL:   "git@github.com:Scalingo/go-gitremote.git",
		Fetch: "+refs/heads/*:refs/remotes/origin/*",
	}

	if !reflect.DeepEqual(*remote, *remotes[0]) {
		t.Errorf("Expected \n%#v\ngot\n%#v", *remote, *remotes[0])
	}
}

func TestListFromConfigContent2Remote(t *testing.T) {
	remotes, err := ListFromConfigContent(config2Remote)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(remotes) != 2 {
		t.Fatalf("Expeced 1 remote, got %v", len(remotes))
	}

	expectedRemotes := Remotes{{
		Name:  "origin",
		URL:   "git@github.com:Scalingo/go-gitremote.git",
		Fetch: "+refs/heads/*:refs/remotes/origin/*",
	}, {
		Name:  "upstream",
		URL:   "git@github.com:Soulou/go-gitremote.git",
		Fetch: "+refs/heads/*:refs/remotes/upstream/*",
	}}

	if !reflect.DeepEqual(expectedRemotes, remotes) {
		t.Errorf("Expected \n%#v\ngot\n%#v", expectedRemotes, remotes)
	}
}

func TestRemoteRepository(t *testing.T) {
	r := &Remote{URL: "git@github.com:Soulou/pipomolo.git"}
	if r.Repository() != "Soulou/pipomolo.git" {
		t.Errorf("Expected %s, got %s", "Soulou/pipomolo.git", r.Repository())
	}

	r = &Remote{URL: "git@github.com:pipomolo.git"}
	if r.Repository() != "pipomolo.git" {
		t.Errorf("Expected %s, got %s", "pipomolo.git", r.Repository())
	}

	r = &Remote{URL: "ssh://git@github.com:22/pipomolo.git"}
	if r.Repository() != "pipomolo.git" {
		t.Errorf("Expected %s, got %s", "pipomolo.git", r.Repository())
	}
}
