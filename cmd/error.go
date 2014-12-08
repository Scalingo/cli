package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/users"
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
	newReportError(err).Report()
	rollbar.Wait()
	fmt.Printf("[Error] %v\n", err)
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
		r.FailedRequest = errgo.Cause(err).(*api.RequestFailedError).Req
	}

	return r
}

func newSysinfo() Sysinfo {
	u, _ := user.Current()
	s := Sysinfo{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Username:  u.Username,
		GoVersion: runtime.Version(),
	}
	return s
}
