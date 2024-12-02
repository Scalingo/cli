package logger

import (
	"regexp"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type RedactionOption struct {
	Field       string
	Regexp      *regexp.Regexp
	ReplaceWith string
}

type RedactingFormatter struct {
	logrus.Formatter
	fields []*RedactionOption
}

func (f *RedactingFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if f.Formatter == nil {
		return nil, errors.New("no formatter set")
	}

	for _, redactionOption := range f.fields {
		if redactionOption == nil || redactionOption.Field == "" {
			continue
		}

		// If the field does not exist in this entry, skip it
		_, ok := entry.Data[redactionOption.Field]
		if !ok {
			continue
		}

		replaceWith := "[REDACTED]"
		if redactionOption.ReplaceWith != "" {
			replaceWith = redactionOption.ReplaceWith
		}

		// Replace the whole field if no regexp is provided
		if redactionOption.Regexp == nil {
			entry.Data[redactionOption.Field] = replaceWith
			continue
		}

		// Replace field content according to the regexp
		entry.Data[redactionOption.Field] = redactionOption.Regexp.
			ReplaceAllString(entry.Data[redactionOption.Field].(string), replaceWith)
	}
	return f.Formatter.Format(entry)
}
