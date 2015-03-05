package gitremote

var (
	config1Remote = `
[core]
       repositoryformatversion = 0
       filemode = true
       bare = false
       logallrefupdates = true
[remote "origin"]
       url = git@github.com:Scalingo/go-gitremote.git
       fetch = +refs/heads/*:refs/remotes/origin/*
[branch "master"]
       remote = origin
       merge = refs/heads/master
`
	config2Remote = `
[core]
       repositoryformatversion = 0
       filemode = true
       bare = false
       logallrefupdates = true
[remote "origin"]
       url = git@github.com:Scalingo/go-gitremote.git
       fetch = +refs/heads/*:refs/remotes/origin/*
[branch "master"]
       remote = origin
       merge = refs/heads/master
[remote "upstream"]
       url = git@github.com:Soulou/go-gitremote.git
       fetch = +refs/heads/*:refs/remotes/upstream/*
`
	config2Hosts = `
[core]
       repositoryformatversion = 0
       filemode = true
       bare = false
       logallrefupdates = true
[remote "origin"]
       url = git@github.com:Scalingo/go-gitremote.git
       fetch = +refs/heads/*:refs/remotes/origin/*
[branch "master"]
       remote = origin
       merge = refs/heads/master
[remote "scalingo"]
       url = git@scalingo.com:go-gitremote.git
       fetch = +refs/heads/*:refs/remotes/scalingo/*
`
	configRemotesPrefix = `
[core]
       repositoryformatversion = 0
       filemode = true
       bare = false
       logallrefupdates = true
[remote "origin"]
       url = git@github.com:Scalingo/go-gitremote.git
       fetch = +refs/heads/*:refs/remotes/origin/*
[remote "scalingo-staging"]
       url = git@scalingo.com:go-gitremote-staging.git
       fetch = +refs/heads/*:refs/remotes/scalingo-staging/*
[branch "master"]
       remote = origin
       merge = refs/heads/master
[remote "scalingo"]
       url = git@scalingo.com:go-gitremote.git
       fetch = +refs/heads/*:refs/remotes/scalingo/*
`
)
