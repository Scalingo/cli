package io

import "fmt"

const (
	boldToken      = "\033[1m"
	boldBlueToken  = "\033[1;34m"
	yellowToken    = "\033[33m"
	greenToken     = "\033[32m"
	grayToken      = "\033[90m"
	lightGrayToken = "\033[1;37m"
	resetToken     = "\033[0m"
)

func Bold(s string) string {
	return fmt.Sprintf("%s%s%s", boldToken, s, resetToken)
}

func BoldBlue(s string) string {
	return fmt.Sprintf("%s%s%s", boldBlueToken, s, resetToken)
}

func Green(s string) string {
	return fmt.Sprintf("%s%s%s", greenToken, s, resetToken)
}

func Yellow(s string) string {
	return fmt.Sprintf("%s%s%s", yellowToken, s, resetToken)
}

func Gray(s string) string {
	return fmt.Sprintf("%s%s%s", grayToken, s, resetToken)
}

func LightGray(s string) string {
	return fmt.Sprintf("%s%s%s", lightGrayToken, s, resetToken)
}
