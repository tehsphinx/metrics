package metrics

import (
	client "github.com/influxdata/influxdb1-client"
	"github.com/rcrowley/go-metrics"
)

func newCounter(name string, options ...Option) *counter {
	m := newMetric(name, options...)
	m.suffix = suffCounter

	// was a metric provided? if not create new one.
	mtrx, ok := m.metric.(metrics.Counter)
	if !ok {
		mtrx = metrics.NewCounter()
	}

	t := &counter{
		baseMetric: *m,
		Counter:    mtrx,
		fieldName:  m.name + m.suffix,
	}
	typeCheck := func(m metric) bool {
		_, ok := m.(*counter)
		return ok
	}

	return m.register(t, typeCheck).(*counter)
}

type counter struct {
	metrics.Counter
	baseMetric
	fieldName string
}

// AddPoints adds points to be written to the db.
func (s *counter) AddPoints(pts []client.Point) []client.Point {
	fields := map[string]interface{}{
		s.fieldName: s.Counter.Snapshot().Count(),
	}
	return append(pts, getPoint(s.measurement, fields, s.tags))
}
