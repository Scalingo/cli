package io

import (
	"bytes"
	"strings"
)

func Indent(data string, level int) string {
	indent := genIndent(level)
	lines := strings.Split(data, "\n")
	buffer := new(bytes.Buffer)
	for _, line := range lines {
		buffer.WriteString(indent + line)
	}
	return buffer.String()
}

func genIndent(level int) string {
	buffer := new(bytes.Buffer)
	for i := 0; i < level; i++ {
		buffer.WriteRune(' ')
	}
	return buffer.String()
}
