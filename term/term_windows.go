package term

import (
	"os"
	"strings"
)

// IsTerminal returns false on Windows.
func IsTerminal(f *os.File) bool {
	return false
}

func IsATTY(f *os.File) bool {
	return false
}

// MakeRaw is a no-op on windows. It returns nil.
func MakeRaw(f *os.File) error {
	if strings.HasSuffix(os.Getenv("SHELL"), "/usr/bin/bash") {
		return makeRaw(f)
	}
	return nil
}

// Restore is a no-op on windows. It returns nil.
func Restore(f *os.File) error {
	if strings.HasSuffix(os.Getenv("SHELL"), "/usr/bin/bash") {
		return restore(f)
	}
	return nil
}

// Cols returns 80 on Windows.
func Cols() (int, error) {
	return 80, nil
}

// Lines returns 24 on Windows.
func Lines() (int, error) {
	return 24, nil
}
