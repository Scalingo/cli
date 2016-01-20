package io

import (
	"fmt"
	"io"
	"time"
)

var loadingRunes = "-\\|/"

type Spinner struct {
	stop     chan struct{}
	writer   io.Writer
	PostHook func()
}

func NewSpinner(writer io.Writer) *Spinner {
	return NewSpinnerWithStopChan(writer, make(chan struct{}))
}

func NewSpinnerWithStopChan(writer io.Writer, stop chan struct{}) *Spinner {
	return &Spinner{
		writer: writer,
		stop:   stop,
	}
}

func (s *Spinner) Start() {
	for i := 0; ; i++ {
		select {
		case <-s.stop:
			fmt.Fprint(s.writer, "\b ")
			if s.PostHook != nil {
				s.PostHook()
			}
			return
		default:
		}
		r := loadingRunes[i%len(loadingRunes)]
		time.Sleep(100 * time.Millisecond)
		fmt.Fprintf(s.writer, "\b%c", r)
	}

}

func (s *Spinner) Stop() {
	close(s.stop)
}
