package io

import (
	"fmt"
	"os"
)

func Error(args ...any) {
	fmt.Fprint(os.Stderr, " !     ")
	fmt.Fprintln(os.Stderr, args...)
}

func Errorf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, " !     "+format, args...)
}

func Warning(args ...any) {
	fmt.Print("  /!\\  ")
	fmt.Println(args...)
}

func Warningf(format string, args ...any) {
	fmt.Printf("  /!\\  "+format, args...)
}

func Status(args ...any) {
	fmt.Print("-----> ")
	fmt.Println(args...)
}

func Statusf(format string, args ...any) {
	fmt.Printf("-----> "+format, args...)
}

func Info(args ...any) {
	fmt.Print("       ")
	fmt.Println(args...)
}

func Infof(format string, args ...any) {
	fmt.Printf("       "+format, args...)
}
