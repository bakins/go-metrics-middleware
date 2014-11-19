package mmw

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/bakins/go-metrics-map"
	h "github.com/bakins/test-helpers"
)

func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func TestHandler(t *testing.T) {
	s := mapsink.New()

	c := metrics.DefaultConfig("test")
	c.EnableHostname = false
	c.EnableRuntimeMetrics = false

	m, err := metrics.New(c, s)

	h.Ok(t, err)
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, "Hello World\n")
	})

	mw := New(m)
	w := httptest.NewRecorder()

	mw.Handler(handler, "testing").ServeHTTP(w, newRequest("GET", "/foo"))

	h.Assert(t, w.Body.String() == "Hello World\n", "body does not match")

	v, ok := s.Get("test.testing.count")
	h.Assert(t, ok, "key not found")
	h.Assert(t, v == 1, "value is not 1")

	v, ok = s.Get("test.testing.time")
	h.Assert(t, ok, "key not found")
}

func TestHandlerWrapper(t *testing.T) {
	s := mapsink.New()

	c := metrics.DefaultConfig("test")
	c.EnableHostname = false
	c.EnableRuntimeMetrics = false

	m, err := metrics.New(c, s)

	h.Ok(t, err)
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, "Hello World\n")
	})

	mw := New(m)
	w := httptest.NewRecorder()

	f := mw.HandlerWrapper("testing")
	f(handler).ServeHTTP(w, newRequest("GET", "/foo"))

	h.Assert(t, w.Body.String() == "Hello World\n", "body does not match")

	v, ok := s.Get("test.testing.count")
	h.Assert(t, ok, "key not found")
	h.Assert(t, v == 1, "value is not 1")

	v, ok = s.Get("test.testing.time")
	h.Assert(t, ok, "key not found")
}

func ExampleMiddleware_Handler() {
	// create a sink to use
	inm := metrics.NewInmemSink(10*time.Millisecond, 50*time.Millisecond)

	// create a default metrics config, without hostname
	conf := metrics.DefaultConfig("test")
	conf.EnableHostname = false

	// now use that config and sink with metrics
	m, _ := metrics.New(conf, inm)

	// simple http handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, "Hello World\n")
	})

	// create a new mw wrapping the metrics
	mw := New(m)

	// wrap the handler and add to a route
	http.Handle("/foo", mw.Handler(handler, "testing"))

	// every time "/foo" is access, the middleware will emit two metrics
	// it will increment the counter "test.testing.count"
	// and add a time sample to "test.testing.time"

	http.Handle("/bar", mw.Handler(handler, "another"))

	// every time "/bar" is access, the middleware will emit two metrics
	// it will increment the counter "test.another.count"
	// and add a time sample to "test.another.time"

	http.ListenAndServe(":8080", nil)
}
