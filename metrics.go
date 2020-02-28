// Package metrics is based on github.com/rcrowley/go-metrics and loosely inspired
// by github.com/vrischmann/go-metrics-influxdb. The main difference of this package
// and `vrischmann/go-metrics-influxdb` is the API. The API is based on the `go-metrics`
// api and separates the registry (built into the Reporter) from the measurement/metrics.
// This enables a global Reporter and the creation of independant metrics distributed across
// the codebase without passing the Reporter around. A reporter can however be injected
// into a metrics if needed.
package metrics

import (
	"errors"

	"github.com/rcrowley/go-metrics"
)

// Register registers a custom metric to the default reporter. Data points from a registered
// metric will be collected via the AddPoints endpoint and then sent to influxDB.
//
// This function is only needed for custom metrics. Functions creating a metric
// in this package register themselves.
func Register(name string, metric Metric, options ...Option) error {
	options = append(options, WithMetric(metric))

	switch metric.(type) {
	case metrics.Counter:
		newCounter(name, options...)
	case metrics.Gauge:
		newGauge(name, options...)
	case metrics.GaugeFloat64:
		newGaugeFloat64(name, options...)
	case metrics.Timer:
		newTimer(name, options...)
	case metrics.Meter:
		newMeter(name, options...)
	case metrics.Histogram:
		newHistogram(name, options...)
	default:
		return errors.New("unknown metrics type")
	}

	return nil
}

// NewCounter creates a new counter or retrieves an existing counter with the same name.
func NewCounter(name string, options ...Option) metrics.Counter {
	return newCounter(name, options...)
}

// NewGauge creates a new gauge or retrieves an existing gauge with the same name.
func NewGauge(name string, options ...Option) metrics.Gauge {
	return newGauge(name, options...)
}

// NewGaugeFloat64 creates a new gauge with float64 or retrieves an existing gauge with the same name.
func NewGaugeFloat64(name string, options ...Option) metrics.GaugeFloat64 {
	return newGaugeFloat64(name, options...)
}

// NewTimer creates a new timer or retrieves an existing timer with the same name.
func NewTimer(name string, options ...Option) metrics.Timer {
	return newTimer(name, options...)
}

// NewMeter creates a new meter or retrieves an existing meter with the same name.
func NewMeter(name string, options ...Option) metrics.Meter {
	return newMeter(name, options...)
}

// NewHistogram creates a new histogram or retrieves an existing histogram with the same name.
// By default this creates a uniform sample with a reservoir size of 100.
// Provide a different metric via the WithMetric option:
// e.g. WithMetric(metrics.NewHistogram(metrics.NewUniformSample(100)))
func NewHistogram(name string, options ...Option) metrics.Histogram {
	return newHistogram(name, options...)
}
