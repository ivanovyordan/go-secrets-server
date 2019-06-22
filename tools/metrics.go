package tools

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	callsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "http",
			Subsystem: "request",
			Name:      "endpoint_calls_total",
			Help:      "Total number of API endpoint calls",
		},
		[]string{"method", "path"},
	)

	requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http",
		Subsystem: "request",
		Name:      "endpoint_duration_seconds",
		Help:      "Time (in seconds) spent serving HTTP requests",
		Buckets:   []float64{0.5, 0.9, 0.95},
	}, []string{"method", "route"})
)

func InitMetrics() {
	prometheus.MustRegister(callsTotal)
	prometheus.MustRegister(requestDuration)
}

func MetricsMiddleware(router http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		start := time.Now()

		route := mux.CurrentRoute(request)
		routeName := route.GetName()

		callsTotal.WithLabelValues(request.Method, routeName).Inc()
		router.ServeHTTP(response, request)

		duration := time.Since(start).Seconds()
		requestDuration.WithLabelValues(request.Method, routeName).Observe(float64(duration))
	})
}

func GetMetrics() http.Handler {
	return promhttp.Handler()
}
