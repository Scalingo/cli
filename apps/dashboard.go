package apps

import (
	"fmt"

	"github.com/pkg/browser"
	"gopkg.in/errgo.v1"
)

func Dashboard(appName string, region string) error {
	url := fmt.Sprintf("https://dashboard.scalingo.com/apps/%s/%s", region, appName)

	if err := browser.OpenURL(url); err != nil {
		return errgo.Notef(err, "fail to open dashboard in browser")
	}

	return nil
}
