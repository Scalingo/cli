package debug

import (
	"log"
	"os"
)

var (
	Enable bool
)

func init() {
	if os.Getenv("DEBUG") == "1" {
		Enable = true
	}
}

func Println(vars ...interface{}) {
	if Enable {
		log.Println(vars...)
	}
}

func Printf(format string, vars ...interface{}) {
	if Enable {
		log.Printf(format, vars...)
	}
}
