package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/debug"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/session"
	"github.com/Scalingo/go-scalingo"
	"github.com/Soulou/errgo-rollbar"
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
	rollbar.ErrorWithStack(rollbar.ERR, r.Error, errgorollbar.BuildStack(r.Error), fields...)
}

func errorQuit(err error) {
	if scalingo.IsRequestFailedError(errgo.Cause(err)) {
		code := errgo.Cause(err).(*scalingo.RequestFailedError).Code
		if code == 401 {
			session.DestroyToken()
		}
	}

	newReportError(err).Report()
	rollbar.Wait()
	io.Error("An error occurred:")
	debug.Println(errgo.Details(err))
	fmt.Println(io.Indent(err.Error(), 7))
	os.Exit(1)
}

func newReportError(err error) *ReportError {
	r := &ReportError{
		Time:    time.Now(),
		User:    config.AuthenticatedUser,
		Error:   err,
		Command: strings.Join(os.Args, " "),
		Version: config.Version,
		System:  newSysinfo(),
	}

	if scalingo.IsRequestFailedError(errgo.Cause(err)) {
		r.FailedRequest = errgo.Cause(err).(*scalingo.RequestFailedError).Req.HTTPRequest
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
