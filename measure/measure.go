package measure

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

//
var (
	// FunctionDurations - Create a summary to track elapsed time of our key functions
	FunctionDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "function_durations_seconds",
			Help:       "Function timings for Rancher Exporter",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}, []string{"pkg", "fnc"})

	// FunctionCountTotal - Create a counter to track total executions of the functions
	FunctionCountTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "function_count_totals",
			Help: "total count of function calls",
		}, []string{"pkg", "fnc"})

	start = time.Now()
)

// Init registers the prometheus metrics for the measurement of the exporter itsself.
func Init() {

	prometheus.MustRegister(FunctionDurations)
	prometheus.MustRegister(FunctionCountTotal)

}
