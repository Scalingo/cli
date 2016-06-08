package gitremote

import "testing"

var (
	config1RemoteDeleted = `
[core]
       repositoryformatversion = 0
       filemode = true
       bare = false
       logallrefupdates = true
[branch "master"]
       remote = origin
       merge = refs/heads/master
`
	config2RemoteDeleted = `
[core]
       repositoryformatversion = 0
       filemode = true
       bare = false
       logallrefupdates = true
[branch "master"]
       remote = origin
       merge = refs/heads/master
[remote "upstream"]
       url = git@github.com:Soulou/go-gitremote.git
       fetch = +refs/heads/*:refs/remotes/upstream/*
`
)

func TestDeleteFromContent(t *testing.T) {
	res := deleteFromContent(config1Remote, &Remote{Name: "origin"})
	if res != config1RemoteDeleted {
		t.Errorf("expected:\n%#v\n--- got:\n%#v", config1RemoteDeleted, res)
	}

	// No affect if done twice
	res = deleteFromContent(res, &Remote{Name: "origin"})
	if res != config1RemoteDeleted {
		t.Errorf("expected:\n%#v\n--- got:\n%#v", config1RemoteDeleted, res)
	}

	res = deleteFromContent(config2Remote, &Remote{Name: "origin"})
	if res != config2RemoteDeleted {
		t.Errorf("expected:\n%#v\n--- got:\n%#v", config2RemoteDeleted, res)
	}

	res = deleteFromContent(res, &Remote{Name: "upstream"})
	if res != config1RemoteDeleted {
		t.Errorf("expected:\n%#v\n--- got:\n%#v", config1RemoteDeleted, res)
	}
}
