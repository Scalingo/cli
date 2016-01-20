package io

import (
	"fmt"
	"os"
	"time"
)

var loadingRunes = "-\\|/"

type Spinner struct {
	stop     chan struct{}
	fd       *os.File
	PostHook func()
}

func NewSpinner(fd *os.File) *Spinner {
	return NewSpinnerWithStopChan(fd, make(chan struct{}))
}

func NewSpinnerWithStopChan(fd *os.File, stop chan struct{}) *Spinner {
	return &Spinner{
		fd:   fd,
		stop: stop,
	}
}

func (s *Spinner) Start() {
	for i := 0; ; i++ {
		select {
		case <-s.stop:
			fmt.Fprint(s.fd, "\b ")
			if s.PostHook != nil {
				s.PostHook()
			}
			return
		default:
		}
		r := loadingRunes[i%len(loadingRunes)]
		time.Sleep(100 * time.Millisecond)
		fmt.Fprintf(s.fd, "\b%c", r)
	}

}

func (s *Spinner) Stop() {
	close(s.stop)
}
