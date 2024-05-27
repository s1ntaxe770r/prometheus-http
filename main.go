package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_request_count",
		Help: "Total number of HTTP requests received",
	})

	RequestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests in mili seconds",
		Buckets: prometheus.LinearBuckets(0.001, 0.005, 10),
	})
)

type pingResponse struct {
	Message string `json:"message"`
}

func main() {
	// Seed the random number generator for better randomness

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("handling new get request")
		startTime := time.Now()
		RequestCounter.Inc() // Increment request counter

		// Generate random number between 1 and 10
		randomNumber := rand.Intn(10) + 1

		var response pingResponse
		response.Message = "hello"

		if randomNumber < 5 {
			// Sleep for 3 seconds
			time.Sleep(3 * time.Second)
		}

		elapsed := time.Since(startTime)
		RequestDuration.Observe(elapsed.Seconds()) // Observe request duration

		// Set content type and encode response as JSON
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			fmt.Println("Error encoding response:", err)
			return
		}
		slog.Info("processed request", "duration", elapsed)
	})

	// Register metrics handler
	prometheus.MustRegister(RequestCounter, RequestDuration)

	// Start HTTP server
	slog.Info("Starting server on port 8080")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
