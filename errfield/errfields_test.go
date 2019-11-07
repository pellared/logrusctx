package errfield_test

import (
	"errors"
	"os"

	"github.com/pellared/logrusutil/errfield"

	log "github.com/sirupsen/logrus"
)

func Example() {
	log.SetOutput(os.Stdout)

	// setup the errfield.Formatter
	log.SetFormatter(&errfield.Formatter{
		Formatter: &log.TextFormatter{DisableTimestamp: true},
	})

	// use errfield.Add to add fields
	err := errors.New("something failed")
	err = errfield.Add(err, "foo", "bar")
	err = errfield.Add(err, "fizz", "buzz")
	log.WithError(err).Error("crash")

	// Output: level=error msg=crash error="something failed" fizz=buzz foo=bar
}