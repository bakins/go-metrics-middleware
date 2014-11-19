package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/armon/go-metrics"
	"github.com/bakins/go-metrics-middleware"
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

	mw := mmw.New(m)

	http.Handle("/foo", mw.Handler(handler, "foo"))

	http.Handle("/bar", mw.Handler(handler, "bar"))
	http.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(inm.Data())
	}))

	http.ListenAndServe(":8080", nil)
}
