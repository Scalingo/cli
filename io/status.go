package io

import "fmt"

func Warning(args ...interface{}) {
	fmt.Print("  /!\\  ")
	fmt.Println(args...)
}

func Status(args ...interface{}) {
	fmt.Print("-----> ")
	fmt.Println(args...)
}

func Info(args ...interface{}) {
	fmt.Print("       ")
	fmt.Println(args...)
}
