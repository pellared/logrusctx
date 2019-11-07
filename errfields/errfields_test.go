package errfields_test

import (
	"errors"
	"os"

	"github.com/pellared/logrusctx/errfields"

	log "github.com/sirupsen/logrus"
)

func Example() {
	log.SetOutput(os.Stdout)

	// setup the errfields.Formatter
	log.SetFormatter(&errfields.Formatter{
		Formatter: &log.TextFormatter{DisableTimestamp: true},
	})

	// use errfields.Add to add fields
	err := errors.New("something failed")
	err = errfields.Add(err, "foo", "bar")
	err = errfields.Add(err, "fizz", "buzz")
	log.WithError(err).Error("crash")

	// Output: level=error msg=crash error="something failed" fizz=buzz foo=bar
}
