package logctx_test

import (
	"context"
	"os"
	"sync/atomic"
	"time"

	"github.com/pellared/logrus/logctx"

	log "github.com/sirupsen/logrus"
)

func Example_reqID() {
	log.SetOutput(os.Stdout)
	reqID := "we232s75tyg9rev"                                            // in reality randomly generated
	timestamp, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00") // hardcode the log timestamp

	// setting contextual log entry
	ctx := logctx.New(context.Background(), log.WithField("ReqID", reqID))

	// retrieving context log entry, adding some data and emitting the log
	logctx.From(ctx).WithTime(timestamp).WithField("foo", "bar").Info("foobar created")

	// Output: time="2012-11-01T22:08:41Z" level=info msg="foobar created" ReqID=we232s75tyg9rev foo=bar
}

func Example_goroutineID() {
	log.SetOutput(os.Stdout)
	timestamp, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00") // hardcode the log timestamp

	// setting up GoroutineIDs
	const LogFieldGoroutineID = "grtnID"
	const LogFieldGoroutineParentID = "grtnPrntID"
	var goroutineIDCounter int64

	// set context log entry for the main goroutine
	logEntry := log.WithField(LogFieldGoroutineID, atomic.AddInt64(&goroutineIDCounter, 1)).WithTime(timestamp)
	ctx := logctx.New(context.Background(), logEntry)

	// spawnGoroutine creates runs new goroutine with contextual log entries that has goroutine IDs
	// returns channel which closes when the goroutine
	// logs error if goroutine panicked
	spawnGoroutine := func(ctx context.Context, fn func(context.Context)) <-chan interface{} {
		entry := logctx.From(ctx)
		if gortnID, ok := entry.Data[LogFieldGoroutineID].(int64); ok {
			entry = entry.WithField(LogFieldGoroutineParentID, gortnID)
		}
		entry = entry.WithField(LogFieldGoroutineID, atomic.AddInt64(&goroutineIDCounter, 1))
		newCtx := logctx.New(ctx, entry)
		done := make(chan interface{})
		go func() {
			defer func() {
				if r := recover(); r != nil {
					entry.
						//	WithField("stack", string(debug.Stack())).
						WithField("panic", r).
						Error("goroutine panicked")
				}
				close(done)
			}()
			fn(newCtx)
		}()
		return done
	}

	<-spawnGoroutine(ctx, func(ctx context.Context) {
		logEntry := logctx.From(ctx).WithField("foo", "bar")
		logEntry.Info("first child goroutine started")

		<-spawnGoroutine(ctx, func(ctx context.Context) {
			logctx.From(ctx).WithField("fizz", "buzz").Info("second child goroutine")

			<-spawnGoroutine(ctx, func(ctx context.Context) {
				panic("panic from third child")
			})
		})

		logEntry.Info("first child goroutine finished")
	})

	// Output:
	// time="2012-11-01T22:08:41Z" level=info msg="first child goroutine started" foo=bar grtnID=2 grtnPrntID=1
	// time="2012-11-01T22:08:41Z" level=info msg="second child goroutine" fizz=buzz grtnID=3 grtnPrntID=2
	// time="2012-11-01T22:08:41Z" level=error msg="goroutine panicked" grtnID=4 grtnPrntID=3 panic="panic from third child"
	// time="2012-11-01T22:08:41Z" level=info msg="first child goroutine finished" foo=bar grtnID=2 grtnPrntID=1
}
