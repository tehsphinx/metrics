package main

import (
	"math/rand"
	"time"

	"github.com/tehsphinx/metrics"
)

func main() {
	// Initialize reporter and set it as default.
	rep := metrics.NewReporter("http://localhost:8086", "metrics", metrics.Interval(1*time.Second))
	metrics.SetDefaultReporter(rep)

	// Start reporting all registered or to be registered metrics.
	go rep.Run()

	// Anywhere else in the code:

	// Create and register a new gauge metric.
	m := metrics.NewGauge("Test", metrics.WithMeasurement("measure"))
	for {
		time.Sleep(600 * time.Millisecond)
		n := rand.Int63n(50)
		// Update the gauge metric value.
		m.Update(n)
	}
}
