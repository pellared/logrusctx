package logctx

import (
	"context"

	"github.com/sirupsen/logrus"
)

// DefaultLogger is used to create a new LogEntry
// if there is no LogEntry within the context.
// Set it in your application's Compostion Root.
var DefaultLogger *logrus.Logger = logrus.StandardLogger()

type contextKey struct{}

// New returns a a copy of parent context and adds the provided log entry.
// Used to set the contextual log entry.
func New(ctx context.Context, logEntry *logrus.Entry) context.Context {
	return context.WithValue(ctx, contextKey{}, logEntry)
}

// From returns the contextual log entry.
func From(ctx context.Context) *logrus.Entry {
	if entry, ok := ctx.Value(contextKey{}).(*logrus.Entry); ok {
		return entry
	}
	// handling case when WithLogger was not invoked for given context
	return logrus.NewEntry(DefaultLogger)
}
