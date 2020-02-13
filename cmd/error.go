package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/go-scalingo/debug"
	httpclient "github.com/Scalingo/go-scalingo/http"
	"github.com/Scalingo/go-utils/errors"
	"github.com/stvp/rollbar"
	"gopkg.in/errgo.v1"
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
	fields := []*rollbar.Field{&rollbar.Field{
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

func errorQuit(err error) {
	currentUser, autherr := config.C.CurrentUser()
	if autherr != nil {
		debug.Println("Fail to get current user")
		debug.Println(errgo.Details(err))
	}

	newReportError(currentUser, err).Report()
	rollbar.Wait()

	rootError := errors.ErrgoRoot(err)
	if httpclient.IsRequestFailedError(rootError) {
		requestFailedErr := rootError.(*httpclient.RequestFailedError)
		if requestFailedErr.Code == 401 {
			if currentUser != nil {
				io.Errorf("You are currently logged in as %s.\n", currentUser.Username)
				io.Errorf("Are you sure %s is a collaborator of this app?\n", currentUser.Username)
			} else {
				io.Errorf("Failed to read credentials for current user: %v", autherr)
			}
		} else if requestFailedErr.Code == 404 {
			notFoundErr := requestFailedErr.APIError.(httpclient.NotFoundError)
			if notFoundErr.Resource == "app" {
				// apiURL contains something like:
				// "https://api.agora-fr1.scalingo.com/v1"
				apiURL, _ := url.Parse(requestFailedErr.Req.URL)
				region := strings.Split(apiURL.Host, ".")[1]
				io.Errorf("The application was not found on the region %s.\n", region)
				io.Error("You can try on a different region with 'scalingo --region osc-fr1 ...'.")
				io.Error("")
				io.Error("List of available regions is accessible with 'scalingo regions'.")
			} else {
				io.Error("An error occured:")
				debug.Println(errgo.Details(err))
				fmt.Println(io.Indent(err.Error(), 7))
			}
		}
	} else {
		io.Error("An error occured:")
		debug.Println(errgo.Details(err))
		fmt.Println(io.Indent(err.Error(), 7))
	}

	os.Exit(1)
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

	rootError := errors.ErrgoRoot(err)
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
