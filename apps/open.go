package apps

import (
	"context"
	"fmt"

	"github.com/pkg/browser"

	"github.com/Scalingo/go-utils/errors/v3"
)

func Open(ctx context.Context, appName string, region string) error {
	url := fmt.Sprintf("https://%s.%s.scalingo.io/", appName, region)

	err := browser.OpenURL(url)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to open app in browser")
	}

	return nil
}
