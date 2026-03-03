package io

import "fmt"

func Error(args ...any) {
	fmt.Print(" !     ")
	fmt.Println(args...)
}

func Warning(args ...any) {
	fmt.Print("  /!\\  ")
	fmt.Println(args...)
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
