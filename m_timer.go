package metrics

import (
	"time"

	client "github.com/influxdata/influxdb1-client"
	"github.com/rcrowley/go-metrics"
)

// Timer implements go-metrics.Timer and possibly adds a bit functionality.
type Timer interface {
	metrics.Timer

	TimeThis() func()
}

func newTimer(name string, options ...Option) *timer {
	m := newMetric(name, options...)
	m.suffix = suffTimer

	// was a metric provided? if not create new one.
	mtrx, ok := m.metric.(metrics.Timer)
	if !ok {
		mtrx = metrics.NewTimer()
	}

	t := &timer{
		baseMetric:  *m,
		Timer:       mtrx,
		fieldName:   m.name + m.suffix,
		percentiles: []float64{0.5, 0.75, 0.95, 0.99, 0.999, 0.9999},
		buckets: []string{count, max, mean, min, p50, p75, p95, p99, p999, p9999,
			stddev, variance, m1, m5, m15, meanrate},
	}
	t.bucketTags = buildBucketTags(t.buckets, t.tags)
	t.bucketVals = buildBucketVals(t.buckets, t.fieldName)
	return m.register(t).(*timer)
}

type timer struct {
	metrics.Timer
	baseMetric
	fieldName   string
	percentiles []float64
	buckets     []string
	bucketTags  map[string]map[string]string
	bucketVals  map[string]map[string]interface{}
}

// AddPoints adds points to be written to the db.
func (s *timer) AddPoints(pts []client.Point) []client.Point {
	var (
		ms = s.Timer.Snapshot()
		ps = ms.Percentiles(s.percentiles)
	)

	for _, bucket := range s.buckets {
		fields := s.bucketVals[bucket]
		fields[s.fieldName] = getValue(bucket, ms, ps)

		pts = append(pts, getPoint(s.measurement, fields, s.bucketTags[bucket]))
	}
	return pts
}

func getValue(bucket string, snapshot metrics.Timer, percentiles []float64) float64 {
	var val float64
	switch bucket {
	case count:
		val = float64(snapshot.Count())
	case max:
		val = float64(snapshot.Max())
	case mean:
		val = snapshot.Mean()
	case min:
		val = float64(snapshot.Min())
	case stddev:
		val = snapshot.StdDev()
	case variance:
		val = snapshot.Variance()
	case p50:
		val = percentiles[0]
	case p75:
		val = percentiles[1]
	case p95:
		val = percentiles[2]
	case p99:
		val = percentiles[3]
	case p999:
		val = percentiles[4]
	case p9999:
		val = percentiles[5]
	case m1:
		val = snapshot.Rate1()
	case m5:
		val = snapshot.Rate5()
	case m15:
		val = snapshot.Rate15()
	case meanrate:
		val = snapshot.RateMean()
	}
	return val
}

// TimeThis measure starts a timer and returns a function to stop the time and report it.
// Can be used as `defer timer.TimeThis()()`.
func (s *timer) TimeThis() func() {
	t := time.Now()
	return func() {
		s.UpdateSince(t)
	}
}
