// Package errfield adds possibility to wrap errors with fields and then log them in structured way.
package errfield

import (
	"errors"

	"github.com/sirupsen/logrus"
)

// Formatter decorates logrus.Formatter to add error fields under to the log entry.
// Implements logrus.Formatter.
type Formatter struct {
	logrus.Formatter
}

// Format implements logrus.Formatter.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var e *Error
	err, ok := entry.Data[logrus.ErrorKey].(error)
	if ok && errors.As(err, &e) {
		for key, value := range e.Fields {
			entry.Data[key] = value
		}
	}
	return f.baseFormatter().Format(entry)
}

func (f *Formatter) baseFormatter() logrus.Formatter {
	if f.Formatter == nil {
		return &logrus.TextFormatter{}
	}
	return f.Formatter
}
