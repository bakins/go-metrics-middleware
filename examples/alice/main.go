package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/armon/go-metrics"
	"github.com/bakins/go-metrics-middleware"
	"github.com/bakins/net-http-recover"
	"github.com/gorilla/handlers"
	"github.com/justinas/alice"
)

func main() {
	inm := metrics.NewInmemSink(10*time.Second, 5*time.Minute)

	conf := metrics.DefaultConfig("test")
	conf.EnableHostname = false

	m, _ := metrics.New(conf, inm)

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, "Hello World\n")
	})

	chain := alice.New(
		func(h http.Handler) http.Handler {
			return handlers.CombinedLoggingHandler(os.Stdout, h)
		},
		handlers.CompressHandler,
		func(h http.Handler) http.Handler {
			return recovery.Handler(os.Stderr, h, true)
		})
	mw := mmw.New(m)

	// Use the wrapper which returns a constructor.
	http.Handle("/foo", chain.Append(mw.HandlerWrapper("foo")).Then(handler))

	// Or wrap the handler directly
	http.Handle("/bar", chain.Then(mw.Handler(handler, "bar")))

	http.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(inm.Data())
	}))

	http.ListenAndServe(":8080", nil)
}
