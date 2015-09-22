package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Soulou/errgo-rollbar"
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/stvp/rollbar"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/session"
	"github.com/Scalingo/cli/users"
)

type Sysinfo struct {
	OS        string
	Arch      string
	Username  string
	GoVersion string
}

type ReportError struct {
	Time          time.Time
	User          *users.User
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
	rollbar.ErrorWithStack(rollbar.ERR, r.Error, errgorollbar.BuildStack(r.Error), fields...)
}

func errorQuit(err error) {
	if errgo.Cause(err) == api.ErrLoginAborted {
		fmt.Printf("... %v\n", err)
		os.Exit(1)
	}

	if api.IsRequestFailedError(errgo.Cause(err)) {
		code := errgo.Cause(err).(*api.RequestFailedError).Code
		if code == 401 {
			session.DestroyToken()
		}
	}

	newReportError(err).Report()
	rollbar.Wait()
	io.Error("An error occured:")
	fmt.Println(io.Indent(err.Error(), 7))
	os.Exit(1)
}

func newReportError(err error) *ReportError {
	r := &ReportError{
		Time:    time.Now(),
		User:    api.CurrentUser,
		Error:   err,
		Command: strings.Join(os.Args, " "),
		Version: config.Version,
		System:  newSysinfo(),
	}

	if api.IsRequestFailedError(errgo.Cause(err)) {
		r.FailedRequest = errgo.Cause(err).(*api.RequestFailedError).Req.HTTPRequest
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
