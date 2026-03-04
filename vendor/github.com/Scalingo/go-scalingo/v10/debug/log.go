package debug

import (
	"log"
	"os"
)

var (
	logger = log.New(os.Stderr, "[DEBUG]", log.LstdFlags)
	Enable bool
)

func init() {
	if os.Getenv("DEBUG") == "1" {
		Enable = true
	}
}

func Println(vars ...any) {
	if Enable {
		logger.Println(vars...)
	}
}

func Printf(format string, vars ...any) {
	if Enable {
		logger.Printf(format, vars...)
	}
}
