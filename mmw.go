// package mmw provides a generic http.Handler middleware for go-metrics.
package mmw

import (
	"fmt"
	"net/http"
	"time"

	"github.com/armon/go-metrics"
)

type metricsHandler struct {
	handler  http.Handler
	mw       *Middleware
	timeKey  []string
	countKey []string
}

// Middleware is simple wrapper around go-metrics.
type Middleware struct {
	sink *metrics.Metrics
}

// New creates a new Middleware
func New(sink *metrics.Metrics) *Middleware {
	return &Middleware{sink: sink}
}

// HandlerWrapper wraps Handler and returns an http.Handler. Useful for chains of middleware like
// https://github.com/justinas/alice
func (mw *Middleware) HandlerWrapper(key ...string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return mw.Handler(h, key...)
	}
}

// Handler creates a new metrics handler that implements http.Handler. This wraps
// a handler
func (mw *Middleware) Handler(handler http.Handler, key ...string) *metricsHandler {

	m := &metricsHandler{
		handler: handler,
		mw:      mw,
	}

	m.timeKey = make([]string, len(key)+1)
	copy(m.timeKey, key)
	m.timeKey[len(key)] = "time"

	m.countKey = make([]string, len(key)+1)
	copy(m.countKey, key)
	m.countKey[len(key)] = "count"

	return m
}

// ServeHTTP wraps a handler and records timing and increments a counter
func (m *metricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(m.timeKey)
	now := time.Now()
	m.handler.ServeHTTP(w, r)
	m.mw.sink.AddSample(m.timeKey, float32(time.Since(now).Seconds()*1000))
	m.mw.sink.IncrCounter(m.countKey, 1)
}
