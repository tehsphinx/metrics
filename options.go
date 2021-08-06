package metrics

import (
	"time"

	"github.com/rcrowley/go-metrics"
)

// Option defines an option to be used when creating a new measurement.
type Option func(s *baseMetric)

// WithReporter sets a reporter to use. If not set the global reporter will be used
// that can be set via `metrics.SetDefaultReporter`.
func WithReporter(r Reporter) Option {
	return func(s *baseMetric) {
		s.reporter = r
	}
}

// WithMeasurement provides a ifluxDB measurement name to the metric.
func WithMeasurement(m string) Option {
	return func(s *baseMetric) {
		s.measurement = m
	}
}

// WithTags adds tags to the metric. The tags given to the reporter are being extended and overwritten (if same key).
func WithTags(t map[string]string) Option {
	return func(s *baseMetric) {
		s.tags = t
	}
}

// WithMetric injects a github.com/rcrowley/go-metrics metric instead of creating a new one.
func WithMetric(m interface{}) Option {
	return func(s *baseMetric) {
		s.metric = m
	}
}

// ReporterOption defines an option to be used when creating a reporter.
type ReporterOption func(r *reporter)

// Auth sets user and password for the influxDB connection.
func Auth(user, pass string) ReporterOption {
	return func(r *reporter) {
		r.server.User = user
		r.server.Pass = pass
	}
}

// Interval overwrites the default interval of 10 seconds. The interval is used to send
// the data e.g. every 10 seconds.
func Interval(d time.Duration) ReporterOption {
	return func(r *reporter) {
		r.interval = d
	}
}

// Registry uses the given metric registry instead of the global default registry.
func Registry(reg metrics.Registry) ReporterOption {
	return func(r *reporter) {
		r.registry = reg
	}
}

// Tags allows to add some tags to be added to all metrics sent by this reporter.
func Tags(tags map[string]string) ReporterOption {
	return func(r *reporter) {
		r.tags = tags
	}
}

// Align enables aligning the reported time to the interval.
func Align() ReporterOption {
	return func(r *reporter) {
		r.align = true
	}
}

// WithGCStats enables collection of GC stats for this reporter.
func WithGCStats() ReporterOption {
	return func(r *reporter) {
		captureGCStats(r.registry, r.interval)
	}
}

// WithMemStats enables collection of memory stats for this reporter.
func WithMemStats() ReporterOption {
	return func(r *reporter) {
		captureMemStats(r.registry, r.interval)
	}
}

func withDBClient(client dbClient) ReporterOption {
	return func(r *reporter) {
		r.client = client
	}
}
