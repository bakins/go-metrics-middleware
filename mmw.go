package mmw

import (
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

type Middleware struct {
	sink *metrics.Metrics
}

func New(sink *metrics.Metrics) *Middleware {
	return &Middleware{sink: sink}
}

func (mw *Middleware) HandlerWrapper(key ...string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return mw.Handler(h, key...)
	}
}

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

// ServeHTTP wraps a handler and records timing and increments a timer
func (m *metricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	m.handler.ServeHTTP(w, r)
	m.mw.sink.AddSample(m.timeKey, float32(time.Since(now).Seconds()))
	m.mw.sink.IncrCounter(m.countKey, 1)
}
