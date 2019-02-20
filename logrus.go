package logrus

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

// Middleware used to log the requests with a level based on the HTTP
// response status.
//
// It is also possible to add an error with the LogError method with the request
// context.
//
// The levels are:
// - 2XX => INFO
// - 3XX / 4XX =>  WARNING
// - 5XX => ERROR
type Middleware struct {
	logger *logrus.Logger
}

// NewMiddleware instantiate a new Middleware.
func NewMiddleware(level string, options ...Option) *Middleware {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		panic("failed to parse the log level")
	}

	log := logrus.New()
	log.Level = logLevel

	for _, opt := range options {
		opt(log)
	}

	return &Middleware{
		logger: log,
	}
}

// Wrap a classic http handler.
func (t *Middleware) Wrap(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.ServeHTTP(w, r, handler)
	})
}

// ServeHTTP serve the request then call the next middleware.
func (t *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	entry := logrus.NewEntry(t.logger)
	entry = t.before(entry, r)

	ctx := context.WithValue(r.Context(), EntryKey, entry)

	// Call the next middleware.
	next(w, r.WithContext(ctx))

	latency := time.Since(start)
	res := w.(negroni.ResponseWriter)

	entry = t.after(entry, res, latency)

	statusCode := res.Status()
	if statusCode >= 200 && statusCode < 300 {
		entry.Info()
	} else if statusCode >= 300 && statusCode < 500 {
		entry.Warning()
	} else {
		entry.Error()
	}
}

func (t *Middleware) before(entry *logrus.Entry, req *http.Request) *logrus.Entry {
	return entry.WithFields(logrus.Fields{
		"request-url":    req.URL.String(),
		"request-method": req.Method,
		"user-agent":     req.UserAgent(),
		"x-unique-id":    req.Header.Get("X-Unique-Id"),
	})
}

func (t *Middleware) after(entry *logrus.Entry, res negroni.ResponseWriter, latency time.Duration) *logrus.Entry {
	return entry.WithFields(logrus.Fields{
		"status":     res.Status(),
		"latency_ms": latency.Round(time.Millisecond),
	})
}
