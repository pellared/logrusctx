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
	var goroutineIDCounter int64

	// set context log entry for main goroutine
	logEntry := log.WithField("gortnID", atomic.AddInt64(&goroutineIDCounter, 1))
	ctx := logctx.New(context.Background(), logEntry)

	// spawnGoroutine creats new with contextual log entries that has goroutine IDs
	spawnGoroutine := func(ctx context.Context, fn func(context.Context)) <-chan error {
		entry := logctx.From(ctx)
		if gortnID, ok := entry.Data["gortnID"].(int64); ok {
			entry = entry.WithField("parentGortnID", gortnID)
		}
		entry = entry.WithField("gortnID", atomic.AddInt64(&goroutineIDCounter, 1))
		newCtx := logctx.New(ctx, entry)
		done := make(chan error)
		go func() {
			fn(newCtx)
			close(done)
		}()
		return done
	}

	<-spawnGoroutine(ctx, func(ctx context.Context) {
		logEntry := logctx.From(ctx).WithTime(timestamp).WithField("foo", "bar")
		logEntry.Info("first child goroutine started")

		<-spawnGoroutine(ctx, func(ctx context.Context) {
			logctx.From(ctx).WithTime(timestamp).WithField("bizz", "buzz").Info("second child goroutine")
		})

		logEntry.Info("first child goroutine finished")
	})

	// Output:
	// time="2012-11-01T22:08:41Z" level=info msg="first child goroutine started" foo=bar gortnID=2 parentGortnID=1
	// time="2012-11-01T22:08:41Z" level=info msg="second child goroutine" bizz=buzz gortnID=3 parentGortnID=2
	// time="2012-11-01T22:08:41Z" level=info msg="first child goroutine finished" foo=bar gortnID=2 parentGortnID=1
}
