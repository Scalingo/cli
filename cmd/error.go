package cmd

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/stvp/rollbar"
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-scalingo/v11/debug"
	httpclient "github.com/Scalingo/go-scalingo/v11/http"
	"github.com/Scalingo/go-utils/errors/v3"
)

type Sysinfo struct {
	OS        string
	Arch      string
	Username  string
	GoVersion string
}

type ReportError struct {
	Time          time.Time
	User          *scalingo.User
	Error         error
	Command       string
	Version       string
	System        Sysinfo
	FailedRequest *http.Request
}

func (r *ReportError) Report() {
	fields := []*rollbar.Field{{
		Name: "custom",
		Data: r,
	}}
	if r.User != nil {
		fields = append(fields, &rollbar.Field{
			Name: "person",
			Data: map[string]string{
				"id":       r.User.ID,
				"username": r.User.Username,
				"email":    r.User.Email,
			},
		})
	}
}

func errorQuitWithHelpMessage(ctx context.Context, err error, c *cli.Command, command string) {
	displayError(ctx, err)
	fmt.Print("\n")
	_ = cli.ShowCommandHelp(ctx, c, command)

	os.Exit(1)
}

func errorQuit(ctx context.Context, err error) {
	displayError(ctx, err)

	os.Exit(1)
}

func displayError(ctx context.Context, err error) {
	currentUser, autherr := config.C.CurrentUser(ctx)
	if autherr != nil {
		debug.Println("Fail to get current user")
		debug.Println(err)
	}

	newReportError(currentUser, err).Report()
	rollbar.Wait()

	var requestFailedErr *httpclient.RequestFailedError
	if errors.As(err, &requestFailedErr) {
		displayRequestFailedError(requestFailedErr, currentUser, autherr, err)
		return
	}

	var unknownRegionErr config.UnknownRegionError
	if errors.As(err, &unknownRegionErr) {
		displayScalingoError(unknownRegionErr)
		return
	}

	displayScalingoError(err)
}

func displayScalingoError(err error) {
	io.Error("An error occurred:")
	debug.Println(err)
	message := err.Error()
	fmt.Println(io.Indent(message, 7))
}

func displayRequestFailedError(requestFailedErr *httpclient.RequestFailedError, currentUser *scalingo.User, autherr error, err error) {
	switch requestFailedErr.Code {
	case http.StatusBadRequest: // 400
		// In case of bad request error, we only want to display the API error.
		// The call to strings.ReplaceAll is useful as the API error message may contain a raw \n.
		io.Errorf("%s\n", strings.ReplaceAll(requestFailedErr.Error(), `\n`, "\n"))

	case http.StatusUnauthorized: // 401
		if currentUser != nil {
			io.Errorf("You are currently logged in as %s.\n", currentUser.Username)
			io.Errorf("Are you sure %s is a collaborator of this app with sufficient privileges?\n", currentUser.Username)
		} else {
			io.Errorf("Failed to read credentials for current user: %v", autherr)
		}

	case http.StatusNotFound: // 404
		// we can ignore the returned error as it is already checked earlier in the code
		notFoundErr, _ := requestFailedErr.APIError.(httpclient.NotFoundError)
		if notFoundErr.Resource == "app" {
			// apiURL contains something like:
			// "https://api.agora-fr1.scalingo.com/v1"
			apiURL, _ := url.Parse(requestFailedErr.Req.URL)
			region := strings.Split(apiURL.Host, ".")[1]
			io.Errorf("The application was not found on the region %s.\n", region)
			io.Error("You can try on a different region with 'scalingo --region osc-fr1 ...'.")
			io.Error("")
			io.Error("List of available regions for your account is accessible with 'scalingo regions'.")
		} else {
			io.Error("An error occurred:")
			debug.Println(err)
			fmt.Println(io.Indent(err.Error(), 7))
		}

	case http.StatusUnprocessableEntity:
		unprocessableEntityErr, ok := requestFailedErr.APIError.(httpclient.UnprocessableEntity)
		if !ok {
			// If the error returned by the API has not been correctly parsed by go-scalingo
			displayScalingoError(err)
		} else {
			io.Errorf("422 Unprocessable Content\n%v\n", unprocessableEntityErr.Error())
		}

	default:
		displayScalingoError(err)
	}
}

func newReportError(currentUser *scalingo.User, err error) *ReportError {
	r := &ReportError{
		Time:    time.Now(),
		User:    currentUser,
		Error:   err,
		Command: strings.Join(os.Args, " "),
		Version: config.Version,
		System:  newSysinfo(),
	}

	var requestFailedError *httpclient.RequestFailedError
	if errors.As(err, &requestFailedError) {
		r.FailedRequest = requestFailedError.Req.HTTPRequest
	}

	return r
}

func newSysinfo() Sysinfo {
	var username string
	u, err := user.Current()
	if err != nil {
		username = "n/a"
		rollbar.Error(rollbar.WARN, err)
	} else {
		username = u.Username
	}

	s := Sysinfo{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Username:  username,
		GoVersion: runtime.Version(),
	}
	return s
}
