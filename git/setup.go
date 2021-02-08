package git

import (
	git "github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v4/debug"
)

type SetupParams struct {
	RemoteName     string
	ForcePutRemote bool
}

func Setup(appName string, params SetupParams) error {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return errgo.Notef(err, "fail to initialize the Git repository")
	}

	url, err := getGitEndpoint(appName)
	if err != nil {
		return errgo.Notef(err, "fail to get the Git endpoint of this app")
	}
	debug.Println("Adding Git remote", url)

	err = putRemoteInRepository(repo, params.RemoteName, url, params.ForcePutRemote)
	if err != nil {
		return err
	}

	io.Status("Successfully added the Git remote", params.RemoteName, "on", appName)
	return nil
}

func putRemoteInRepository(repository *git.Repository, remoteName string, url string, force bool) error {
	has, err := repositoryHasRemote(repository, remoteName)
	if err != nil {
		return err
	}
	if has && force {
		err := deleteThenCreateRemoteInRepository(repository, remoteName, url)
		if err != nil {
			return err
		}
		return nil
	}
	err = createRemoteInRepository(repository, remoteName, url)
	if err != nil {
		return err
	}
	return nil
}

func createRemoteInRepository(repository *git.Repository, remoteName string, url string) error {
	_, err := repository.CreateRemote(&gitconfig.RemoteConfig{
		Name: remoteName,
		URLs: []string{url},
	})
	if err != nil {
		errWrapped := utils.WrapError(err, "fail to create the Git remote")
		if err == git.ErrRemoteExists {
			message := "Fail to configure git repository, '" + remoteName + "' remote already exists (use --force option to override)"
			errWrapped = utils.WrapError(errWrapped, message)
		}
		return errWrapped
	}
	return nil
}

func deleteThenCreateRemoteInRepository(repository *git.Repository, remoteName string, url string) error {
	err := repository.DeleteRemote(remoteName)
	if err != nil {
		return utils.WrapError(err, "fail to delete the Git remote")
	}
	return createRemoteInRepository(repository, remoteName, url)
}

func repositoryHasRemote(repository *git.Repository, remoteName string) (bool, error) {
	config, err := repository.Storer.Config()
	if err != nil {
		return false, err
	}

	if _, has := config.Remotes[remoteName]; has {
		return true, nil
	}
	return false, nil
}
