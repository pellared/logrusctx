package errfields

import (
	"errors"

	"github.com/sirupsen/logrus"
)

// LogFormatter decorates logrus.Formatter to add error fields under to the log entry.
// Implements logrus.Formatter.
type LogFormatter struct {
	logrus.Formatter
}

// Format implements logrus.Formatter.
func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var e *Error
	err, ok := entry.Data[logrus.ErrorKey].(error)
	if ok && errors.As(err, &e) {
		for key, value := range e.Fields {
			entry.Data[key] = value
		}
	}
	return f.baseFormatter().Format(entry)
}

func (f *LogFormatter) baseFormatter() logrus.Formatter {
	if f.Formatter == nil {
		return &logrus.TextFormatter{}
	}
	return f.Formatter
}
