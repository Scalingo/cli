package io

import (
	"bytes"
	"strings"
)

func Indent(data string, level int) string {
	indent := genIndent(level)
	lines := strings.Split(data, "\n")
	buffer := new(bytes.Buffer)
	for i, line := range lines {
		if i == len(lines)-1 && line == "" {
			continue
		}
		if i == len(lines)-1 {
			buffer.WriteString(indent + line)
		} else {
			buffer.WriteString(indent + line + "\n")
		}
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
