package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Version = "master"

var info = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "example",
		Name:      "info",
		Help:      "Build and runtime info",
	},
	[]string{"version", "started_at"},
)

var counter_requests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "example",
		Name:      "requests",
		Help:      "Number of requests",
	},
	[]string{"method", "path", "status_code"},
)

var queue = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Namespace: "example",
		Name:      "queue",
		Help:      "Nuber of events in queue",
	},
)

var queueInt int

func main() {
	prometheus.MustRegister(info)
	prometheus.MustRegister(counter_requests)
	prometheus.MustRegister(queue)

	http.Handle("/metrics", promhttp.Handler())

	info.WithLabelValues(Version, time.Now().Format(time.RFC3339)).Set(1)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		counter_requests.WithLabelValues(r.Method, r.URL.Path, "200").Inc()
		fmt.Fprintf(w, "Hello World!\n")
	})

	http.HandleFunc("/inc", func(w http.ResponseWriter, r *http.Request) {
		counter_requests.WithLabelValues(r.Method, r.URL.Path, "200").Inc()
		queue.Inc()
		queueInt++
		fmt.Fprintf(w, "+1 (%d)\n", queueInt)
	})

	http.HandleFunc("/dec", func(w http.ResponseWriter, r *http.Request) {
		if queueInt == 0 {
			counter_requests.WithLabelValues(r.Method, r.URL.Path, "500").Inc()
			w.WriteHeader(500)
			fmt.Fprintf(w, "error: queue is empty\n")
			return
		}

		counter_requests.WithLabelValues(r.Method, r.URL.Path, "200").Inc()
		queue.Dec()
		queueInt--
		fmt.Fprintf(w, "-1 (%d)\n", queueInt)
	})

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		counter_requests.WithLabelValues(r.Method, r.URL.Path, "404").Inc()
		w.WriteHeader(404)
		fmt.Fprintf(w, "404 Not Found\n")
	})

	fmt.Println("Listen on 0.0.0.0:8000, see http://127.0.0.1:8000")
	http.ListenAndServe(":8000", nil)
}
