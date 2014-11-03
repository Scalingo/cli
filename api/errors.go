package api

import (
	"bytes"
	"fmt"
	"strings"
)

type InternalError struct {
	Error string `json:"error"`
}

type UnprocessableEntity struct {
	Errors map[string][]string `json:"errors"`
}

func (err *UnprocessableEntity) Error() string {
	b := new(bytes.Buffer)
	for attr, attrErrs := range err.Errors {
		b.WriteString(fmt.Sprintf("%s → %s\n", attr, strings.Join(attrErrs, ", ")))
	}
	return b.String()
}
