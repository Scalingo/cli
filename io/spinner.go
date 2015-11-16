package io

import (
	"fmt"
	"time"
)

var loadingRunes = "-\\|/"

func Spinner(stop chan struct{}) {
	SpinnerWithPosthook(stop, func() {})
}

func SpinnerWithPosthook(stop chan struct{}, hook func()) {
	for i := 0; ; i++ {
		select {
		case <-stop:
			fmt.Print("\b ")
			hook()
			return
		default:
		}
		r := loadingRunes[i%len(loadingRunes)]
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("\b%c", r)
	}
}
