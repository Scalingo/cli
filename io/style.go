package io

import "fmt"

const (
	boldToken  = "\033[1m"
	resetToken = "\033[0m"
)

func Bold(s string) string {
	return fmt.Sprintf("%s%s%s", boldToken, s, resetToken)
}
