package apps

import (
	"fmt"

	"github.com/pkg/browser"
	"gopkg.in/errgo.v1"
)

func Open(appName string, region string) error {
	url := fmt.Sprintf("https://%s.%s.scalingo.io/", appName, region)

	if err := browser.OpenURL(url); err != nil {
		return errgo.Notef(err, "fail to open app in browser")
	}

	return nil
}