package mmw

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
}
