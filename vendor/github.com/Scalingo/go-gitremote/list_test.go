package gitremote

import (
	"reflect"
	"testing"
)

func TestRemotesHosts(t *testing.T) {
	remotes, err := ListFromConfigContent(config2Remote)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	hosts := []string{"github.com", "github.com"}
	if !reflect.DeepEqual(hosts, remotes.Hosts()) {
		t.Errorf("Expected %#v, got %#v", hosts, remotes.Hosts())
	}
}

func TestRemotesFilterByHost(t *testing.T) {
	remotes, err := ListFromConfigContent(config2Hosts)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	filteredRemotes := remotes.FilterByHost("scalingo.com")
	if !reflect.DeepEqual(Remotes{remotes[1]}, filteredRemotes) {
		t.Errorf("Expected %#v, got %#v", Remotes{remotes[1]}, filteredRemotes)
	}
}

func TestRemotesFilterByName(t *testing.T) {
	remotes, err := ListFromConfigContent(config2Hosts)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	filteredRemotes := remotes.FilterByName("origin")
	if !reflect.DeepEqual(Remotes{remotes[0]}, filteredRemotes) {
		t.Errorf("Expected %#v, got %#v", Remotes{remotes[0]}, filteredRemotes)
	}
}

func TestRemotesFilterByNamePrefix(t *testing.T) {
	remotes, err := ListFromConfigContent(configRemotesPrefix)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	filteredRemotes := remotes.FilterByNamePrefix("scalingo")
	if !reflect.DeepEqual(Remotes{remotes[1], remotes[2]}, filteredRemotes) {
		t.Errorf("Expected %#v, got %#v", Remotes{remotes[1], remotes[2]}, filteredRemotes)
	}
}
