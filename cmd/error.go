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
	"github.com/urfave/cli/v2"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v8"
	"github.com/Scalingo/go-scalingo/v8/debug"
	httpclient "github.com/Scalingo/go-scalingo/v8/http"
	"github.com/Scalingo/go-utils/errors/v2"
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
	//rollbar.ErrorWithStack(rollbar.ERR, r.Error, errgorollbar.BuildStack(r.Error), fields...)
}

func errorQuitWithHelpMessage(err error, ctxCli *cli.Context, command string) {
	displayError(ctxCli.Context, err)
	fmt.Print("\n")
	_ = cli.ShowCommandHelp(ctxCli, command)

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
		debug.Println(errgo.Details(err))
	}

	newReportError(currentUser, err).Report()
	rollbar.Wait()

	rootError := errors.RootCause(err)
	if httpclient.IsRequestFailedError(rootError) {
		displayRequestFailedError(rootError, currentUser, autherr, err)
	} else if _, ok := rootError.(config.UnknownRegionError); ok {
		displayScalingoError(rootError)
	} else {
		displayScalingoError(err)
	}
}

func displayScalingoError(err error) {
	io.Error("An error occurred:")
	debug.Println(errgo.Details(err))
	message := err.Error()
	fmt.Println(io.Indent(message, 7))
}

func displayRequestFailedError(rootError error, currentUser *scalingo.User, autherr error, err error) {
	// we can ignore the returned error as it is already checked earlier in the code
	requestFailedErr, _ := rootError.(*httpclient.RequestFailedError)
	if requestFailedErr.Code == 400 {
		// In case of bad request error, we only want to display the API error.
		// The call to strings.ReplaceAll is useful as the API error message may contain a raw \n.
		io.Errorf("%s\n", strings.ReplaceAll(requestFailedErr.Error(), `\n`, "\n"))
		return
	}
	if requestFailedErr.Code == 401 {
		if currentUser != nil {
			io.Errorf("You are currently logged in as %s.\n", currentUser.Username)
			io.Errorf("Are you sure %s is a collaborator of this app?\n", currentUser.Username)
		} else {
			io.Errorf("Failed to read credentials for current user: %v", autherr)
		}
		return
	}
	if requestFailedErr.Code == 404 {
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
			debug.Println(errgo.Details(err))
			fmt.Println(io.Indent(err.Error(), 7))
		}
		return
	}
	displayScalingoError(err)
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

	rootError := errors.RootCause(err)
	if httpclient.IsRequestFailedError(rootError) {
		r.FailedRequest = rootError.(*httpclient.RequestFailedError).Req.HTTPRequest
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
