package main

import (
	"math/rand"
	"time"

	goMetrics "github.com/rcrowley/go-metrics"
	"github.com/tehsphinx/metrics"
)

func main() {
	// Initialize reporters. Make sure to provide a registry so they don't use the same (default registry).
	rep1 := metrics.NewReporter("http://localhost:8086", "metrics",
		metrics.Interval(1*time.Second), metrics.Registry(goMetrics.NewRegistry()))
	rep2 := metrics.NewReporter("http://localhost:8086", "metrics",
		metrics.Interval(3*time.Second), metrics.Registry(goMetrics.NewRegistry()))

	// Start reporting all registered or to be registered metrics.
	go rep1.Run()
	go rep2.Run()

	go newGauge("metric1", rep1)
	go newGauge("metric2", rep2)

	select {}
}

func newGauge(name string, rep metrics.Reporter) {
	// Create and register a new gauge metric. Provide the reporter to which to register.
	m := metrics.NewGauge(name, metrics.WithMeasurement("measure"), metrics.WithReporter(rep))
	for {
		time.Sleep(500 * time.Millisecond)
		n := rand.Int63n(50)
		// Update the gauge metric value.
		m.Update(n)
	}
}
