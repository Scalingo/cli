package apps

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

func Create(appName string, remote string, buildpack string) error {
	c := config.ScalingoClient()
	app, err := c.AppsCreate(scalingo.AppsCreateOpts{Name: appName})
	if err != nil {
		reqestFailedError, ok := errgo.Cause(err).(*scalingo.RequestFailedError)
		if !ok || reqestFailedError.Code != 402 {
			return errgo.Mask(err, errgo.Any)
		}
		paymentRequiredErr, ok := reqestFailedError.APIError.(scalingo.PaymentRequiredError)
		if !ok || paymentRequiredErr.Name != "free-trial-exceeded" {
			return errgo.Mask(err, errgo.Any)
		}
		// If error is Payment Required and user tries to exceed its free trial
		fmt.Println("You are still in your free trial. If you continue, your free trial will end and you will be billed for your usage of the platform. Do you agree? [Y/n]")
		in, err := readCharFromStdin()
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
		if in != "" && strings.ToUpper(in) != "Y" {
			fmt.Println("Doing nothing")
			return nil
		}

		_, err = c.UpdateUser(scalingo.UpdateUserParams{StopFreeTrial: true})
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}

		app, err = c.AppsCreate(scalingo.AppsCreateOpts{Name: appName})
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
	}

	if buildpack != "" {
		fmt.Println("Installing custom buildpack...")
		_, _, err := c.VariableSet(app.Name, "BUILDPACK_URL", buildpack)
		if err != nil {
			fmt.Println("Failed to set custom buildpack. Please add BUILDPACK_URL=" + buildpack + " to your application environment")
		}
	}

	fmt.Printf("App '%s' has been created\n", app.Name)
	if _, ok := appdetect.DetectGit(); ok && appdetect.AddRemote(app.GitUrl, remote) == nil {
		fmt.Printf("Git repository detected: remote %s added\n→ 'git push %s master' to deploy your app\n", remote, remote)
	} else {
		fmt.Printf("To deploy your application, run these commands in your GIT repository:\n→ git remote add %s %s\n→ git push %s master\n", remote, app.GitUrl, remote)
	}
	return nil
}

// Read a single character on stdin. The string is trimmed of white space.
// If the string is then empty, its value is "Y"
func readCharFromStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return "", nil
	}
	return string(input[0]), nil
}
