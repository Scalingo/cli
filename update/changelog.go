package update

import (
	"context"
	"fmt"
	"strings"

	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/services/github"
	"github.com/Scalingo/go-utils/errors/v3"
)

func ShowLastChangelog() error {
	cliLastRelease, err := github.NewClient().GetLatestRelease(context.Background())
	if err != nil {
		return errors.Wrapf(context.Background(), err, "fail to get last CLI release")
	}

	if cliLastRelease.GetBody() == "" {
		io.Warning("GitHub last release is empty")
		return nil
	}

	fmt.Printf("Changelog of the version %v\n\n", cliLastRelease.GetTagName())
	fmt.Printf("%v\n\n", strings.ReplaceAll(cliLastRelease.GetBody(), "\\r\\n", "\r\n"))

	return nil
}
