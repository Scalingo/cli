package io

import (
	"fmt"

	"github.com/fatih/color"
)

func Error(args ...interface{}) {
	fmt.Print(" !     ")
	fmt.Println(args...)
}

func Errorf(format string, args ...interface{}) {
	fmt.Printf(" !     "+format, args...)
}

func Warning(args ...interface{}) {
	fmt.Print("  /!\\  ")
	fmt.Println(args...)
}

func Status(args ...interface{}) {
	fmt.Print("-----> ")
	fmt.Println(args...)
}

func Statusf(format string, args ...interface{}) {
	fmt.Printf("-----> "+format, args...)
}

func Statusfred(format string, args ...interface{}) {
	color.Set(color.FgRed, color.Bold)
	fmt.Printf("-----> "+format, args...)
	color.Unset()
}

func Info(args ...interface{}) {
	fmt.Print("       ")
	fmt.Println(args...)
}

func Infof(format string, args ...interface{}) {
	fmt.Printf("       "+format, args...)
}
func Infofred(format string, args ...interface{}) {
	color.Set(color.FgRed, color.Bold)
	fmt.Printf("       "+format, args...)
	color.Unset()
}
