package metrics

import (
	client "github.com/influxdata/influxdb1-client"
	"github.com/rcrowley/go-metrics"
)

// GaugeFloat64 implements go-metrics.GaugeFloat64 and possibly adds a bit functionality.
type GaugeFloat64 interface {
	metrics.GaugeFloat64
}

func newGaugeFloat64(name string, options ...Option) *gaugeFloat64 {
	m := newMetric(name, options...)
	m.suffix = suffGauge

	// was a metric provided? if not create new one.
	mtrx, ok := m.metric.(metrics.GaugeFloat64)
	if !ok {
		mtrx = metrics.NewGaugeFloat64()
	}

	t := &gaugeFloat64{
		baseMetric:   *m,
		GaugeFloat64: mtrx,
		fieldName:    m.name + m.suffix,
	}
	typeCheck := func(m metric) bool {
		_, ok := m.(*gaugeFloat64)
		return ok
	}

	return m.register(t, typeCheck).(*gaugeFloat64)
}

type gaugeFloat64 struct {
	metrics.GaugeFloat64
	baseMetric
	fieldName string
}

// AddPoints adds points to be written to the db.
func (s *gaugeFloat64) AddPoints(pts []client.Point) []client.Point {
	fields := map[string]interface{}{
		s.fieldName: s.GaugeFloat64.Snapshot().Value(),
	}
	return append(pts, getPoint(s.measurement, fields, s.tags))
}
