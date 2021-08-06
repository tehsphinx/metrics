package metrics

import (
	client "github.com/influxdata/influxdb1-client"
	"github.com/rcrowley/go-metrics"
)

// Meter implements go-metrics.Meter and possibly adds a bit functionality.
type Meter interface {
	metrics.Meter
}

func newMeter(name string, options ...Option) *meter {
	m := newMetric(name, options...)
	m.suffix = suffMeter

	// was a metric provided? if not create new one.
	mtrx, ok := m.metric.(metrics.Meter)
	if !ok {
		mtrx = metrics.NewMeter()
	}

	t := &meter{
		baseMetric: *m,
		Meter:      mtrx,
		fieldName:  m.name + m.suffix,
		buckets:    []string{count, m1, m5, m15, mean},
	}
	t.bucketTags = buildBucketTags(t.buckets, t.tags)
	t.bucketVals = buildBucketVals(t.buckets, t.fieldName)
	return m.register(t).(*meter)
}

type meter struct {
	metrics.Meter
	baseMetric
	fieldName  string
	buckets    []string
	bucketTags map[string]map[string]string
	bucketVals map[string]map[string]interface{}
}

// AddPoints adds points to be written to the db.
func (s *meter) AddPoints(pts []client.Point) []client.Point {
	ms := s.Meter.Snapshot()

	for _, bucket := range s.buckets {
		var val float64

		switch bucket {
		case count:
			val = float64(ms.Count())
		case m1:
			val = ms.Rate1()
		case m5:
			val = ms.Rate5()
		case m15:
			val = ms.Rate15()
		case mean:
			val = ms.RateMean()
		}

		fields := s.bucketVals[bucket]
		fields[s.fieldName] = val

		point := getPoint(s.measurement, fields, s.bucketTags[bucket])
		pts = append(pts, point)
	}
	return pts
}
