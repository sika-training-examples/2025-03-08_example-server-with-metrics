package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"example/version"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var info = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "example",
		Name:      "info",
		Help:      "Build and runtime info",
	},
	[]string{"version", "started_at", "hostname"},
)

var counter_requests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "example",
		Name:      "requests",
		Help:      "Number of requests",
	},
	[]string{"method", "path", "status_code"},
)

var request_duration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "example",
		Name:      "request_duration_seconds",
		Help:      "Request duration in seconds",
		Buckets:   []float64{.01, .05, .1, .2, .5, 1, 2, 5},
	},
	[]string{"status_code", "path", "method"},
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
	prometheus.MustRegister(request_duration)
	prometheus.MustRegister(queue)

	http.Handle("/metrics", promhttp.Handler())

	var hostname, _ = os.Hostname()

	info.WithLabelValues(version.Version, time.Now().Format(time.RFC3339), hostname).Set(1)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		randomSleep()
		counter_requests.WithLabelValues(r.Method, r.URL.Path, "200").Inc()
		fmt.Fprintf(w, "Hello World!\n")
		request_duration.WithLabelValues(r.Method, r.URL.Path, "200").Observe(time.Since(started).Seconds())
	})

	http.HandleFunc("/inc", func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		randomSleep()
		counter_requests.WithLabelValues(r.Method, r.URL.Path, "200").Inc()
		queue.Inc()
		queueInt++
		fmt.Fprintf(w, "+1 (%d)\n", queueInt)
		request_duration.WithLabelValues(r.Method, r.URL.Path, "200").Observe(time.Since(started).Seconds())
	})

	http.HandleFunc("/dec", func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		randomSleep()
		if queueInt == 0 {
			counter_requests.WithLabelValues(r.Method, r.URL.Path, "500").Inc()
			w.WriteHeader(500)
			fmt.Fprintf(w, "error: queue is empty\n")
			request_duration.WithLabelValues(r.Method, r.URL.Path, "500").Observe(time.Since(started).Seconds())
			return
		}

		counter_requests.WithLabelValues(r.Method, r.URL.Path, "200").Inc()
		queue.Dec()
		queueInt--
		fmt.Fprintf(w, "-1 (%d)\n", queueInt)
		request_duration.WithLabelValues(r.Method, r.URL.Path, "200").Observe(time.Since(started).Seconds())
	})

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		counter_requests.WithLabelValues(r.Method, r.URL.Path, "404").Inc()
		w.WriteHeader(404)
		fmt.Fprintf(w, "404 Not Found\n")
		request_duration.WithLabelValues(r.Method, r.URL.Path, "404").Observe(time.Since(started).Seconds())
	})

	fmt.Printf("Version %s (%s), listen on 0.0.0.0:8000, see http://127.0.0.1:8000\n", version.Version, hostname)
	http.ListenAndServe(":8000", nil)
}

func randomSleep() {
	r := rand.Intn(100)
	if r < 50 {
		time.Sleep(50 * time.Millisecond)
		return
	}
	if r < 90 {
		time.Sleep(150 * time.Millisecond)
		return
	}
	if r < 96 {
		time.Sleep(20 * time.Millisecond)
		return
	}
	if r < 98 {
		time.Sleep(5 * time.Millisecond)
		return
	}
	if r < 99 {
		time.Sleep(250 * time.Millisecond)
		return
	}
	time.Sleep(3000 * time.Millisecond)
}
