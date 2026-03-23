package git

import (
	"context"

	git "github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"

	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v11/debug"
	"github.com/Scalingo/go-utils/errors/v3"
)

type SetupParams struct {
	RemoteName     string
	ForcePutRemote bool
}

func Setup(ctx context.Context, appName string, params SetupParams) error {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to initialize the Git repository")
	}

	url, err := getGitEndpoint(ctx, appName)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get the Git endpoint of this app")
	}
	debug.Println("Adding Git remote", url)

	err = putRemoteInRepository(ctx, repo, params.RemoteName, url, params.ForcePutRemote)
	if err != nil {
		return errors.Wrapf(ctx, err, "configure git remote %s", params.RemoteName)
	}

	io.Status("Successfully added the Git remote", params.RemoteName, "on", appName)
	return nil
}

func putRemoteInRepository(ctx context.Context, repository *git.Repository, remoteName string, url string, force bool) error {
	has, err := repositoryHasRemote(ctx, repository, remoteName)
	if err != nil {
		return errors.Wrapf(ctx, err, "check if git remote %s exists", remoteName)
	}
	if has && force {
		err := deleteThenCreateRemoteInRepository(ctx, repository, remoteName, url)
		if err != nil {
			return errors.Wrapf(ctx, err, "replace git remote %s", remoteName)
		}
		return nil
	}
	err = createRemoteInRepository(ctx, repository, remoteName, url)
	if err != nil {
		return errors.Wrapf(ctx, err, "create git remote %s", remoteName)
	}
	return nil
}

func createRemoteInRepository(ctx context.Context, repository *git.Repository, remoteName string, url string) error {
	_, err := repository.CreateRemote(&gitconfig.RemoteConfig{
		Name: remoteName,
		URLs: []string{url},
	})
	if err != nil {
		errWrapped := errors.Wrapf(ctx, err, "create the Git remote")
		if err == git.ErrRemoteExists {
			message := "Fail to configure git repository, '" + remoteName + "' remote already exists (use --force option to override)"
			errWrapped = errors.Wrap(ctx, errWrapped, message)
		}
		return errWrapped
	}
	return nil
}

func deleteThenCreateRemoteInRepository(ctx context.Context, repository *git.Repository, remoteName string, url string) error {
	err := repository.DeleteRemote(remoteName)
	if err != nil {
		return errors.Wrap(ctx, err, "delete the Git remote")
	}
	return createRemoteInRepository(ctx, repository, remoteName, url)
}

func repositoryHasRemote(ctx context.Context, repository *git.Repository, remoteName string) (bool, error) {
	config, err := repository.Storer.Config()
	if err != nil {
		return false, errors.Wrapf(ctx, err, "fail to get Git repository config for remote %s", remoteName)
	}

	if _, has := config.Remotes[remoteName]; has {
		return true, nil
	}
	return false, nil
}
