package metrics

import (
	client "github.com/influxdata/influxdb1-client"
	"github.com/rcrowley/go-metrics"
)

// Gauge implements go-metrics.Gauge and possibly adds a bit functionality.
type Gauge interface {
	metrics.Gauge
}

func newGauge(name string, options ...Option) *gauge {
	m := newMetric(name, options...)
	m.suffix = suffGauge

	// was a metric provided? if not create new one.
	mtrx, ok := m.metric.(metrics.Gauge)
	if !ok {
		mtrx = metrics.NewGauge()
	}

	t := &gauge{
		baseMetric: *m,
		Gauge:      mtrx,
		fieldName:  m.name + m.suffix,
	}
	return m.register(t).(*gauge)
}

type gauge struct {
	metrics.Gauge
	baseMetric
	fieldName string
}

// AddPoints adds points to be written to the db.
func (s *gauge) AddPoints(pts []client.Point) []client.Point {
	fields := map[string]interface{}{
		s.fieldName: s.Gauge.Snapshot().Value(),
	}
	return append(pts, getPoint(s.measurement, fields, s.tags))
}
