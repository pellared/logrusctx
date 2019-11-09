# logrusutil :hammer: [![GoDoc](https://godoc.org/github.com/pellared/logrusutil?status.svg)](https://godoc.org/github.com/pellared/logrusutil) ![Build Status](https://github.com/pellared/logrusutil/workflows/build/badge.svg) [![golangci](https://golangci.com/badges/github.com/pellared/logrusutil.svg)](https://golangci.com/r/github.com/pellared/logrusutil)

Small, easy to use, yet powerful utility packages for <https://github.com/sirupsen/logrus>.

## `logctx` package [![GoDoc](https://godoc.org/github.com/pellared/logrusutil/logctx?status.svg)](https://godoc.org/github.com/pellared/logrusutil/logctx)

Add a log entry to the context using `logctx.New(ctx, logEntry)`. Retrieve the log entry using `logctx.From(ctx)`.

## `errfield` package [![GoDoc](https://godoc.org/github.com/pellared/logrusutil/errfield?status.svg)](https://godoc.org/github.com/pellared/logrusutil/errfield)

Wrap an error with fields using `errfield.Add(err, "key", value)`. Use `errfield.Formatter` to log the fields in a structured way.
