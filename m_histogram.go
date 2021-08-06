package metrics

import (
	client "github.com/influxdata/influxdb1-client"
	"github.com/rcrowley/go-metrics"
)

// Histogram implements go-metrics.Histogram and possibly adds a bit functionality.
type Histogram interface {
	metrics.Histogram
}

func newHistogram(name string, options ...Option) *histogram {
	m := newMetric(name, options...)
	m.suffix = suffHistogram

	// was a metric provided? if not create new one.
	mtrx, ok := m.metric.(metrics.Histogram)
	if !ok {
		mtrx = metrics.NewHistogram(metrics.NewUniformSample(100))
	}

	t := &histogram{
		baseMetric:  *m,
		Histogram:   mtrx,
		fieldName:   m.name + m.suffix,
		percentiles: []float64{0.5, 0.75, 0.95, 0.99, 0.999, 0.9999},
		buckets: []string{count, max, mean, min, p50, p75, p95, p99, p999, p9999,
			stddev, variance},
	}
	t.bucketTags = buildBucketTags(t.buckets, t.tags)
	t.bucketVals = buildBucketVals(t.buckets, t.fieldName)
	return m.register(t).(*histogram)
}

type histogram struct {
	metrics.Histogram
	baseMetric
	fieldName   string
	percentiles []float64
	buckets     []string
	bucketTags  map[string]map[string]string
	bucketVals  map[string]map[string]interface{}
}

// AddPoints adds points to be written to the db.
func (s *histogram) AddPoints(pts []client.Point) []client.Point {
	var (
		ms  = s.Histogram.Snapshot()
		pct = ms.Percentiles(s.percentiles)
	)

	for _, bucket := range s.buckets {
		var val float64
		switch bucket {
		case count:
			val = float64(ms.Count())
		case max:
			val = float64(ms.Max())
		case mean:
			val = ms.Mean()
		case min:
			val = float64(ms.Min())
		case stddev:
			val = ms.StdDev()
		case variance:
			val = ms.Variance()
		case p50:
			val = pct[0]
		case p75:
			val = pct[1]
		case p95:
			val = pct[2]
		case p99:
			val = pct[3]
		case p999:
			val = pct[4]
		case p9999:
			val = pct[5]
		}

		fields := s.bucketVals[bucket]
		fields[s.fieldName] = val

		pts = append(pts, getPoint(s.measurement, fields, s.bucketTags[bucket]))
	}
	return pts
}
