package metrics

import (
	"testing"
	"time"

	client "github.com/influxdata/influxdb1-client"
	"github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"
	"github.com/tehsphinx/concurrent"
)

func TestNewCounter(t *testing.T) {
	var count = concurrent.NewInt()

	reporter := NewReporter("", "testDB",
		Interval(30*time.Millisecond),
		Tags(map[string]string{
			"reporterTag1": "valRep1",
			"reporterTag2": "valRep2",
			"reporterTag3": "valRep3",
		}),
		Registry(metrics.NewRegistry()),
		withDBClient(&testClient{
			pingCall: func() (duration time.Duration, s string, err error) {
				return 0, "", nil
			},
			writeCall: func(points client.BatchPoints) (response *client.Response, err error) {
				count.Increase()
				assert.Equal(t, "testDB", points.Database)
				assert.True(t, points.Time.Before(time.Now()))
				assert.True(t, points.Time.After(time.Now().Add(-time.Second)))
				assert.Equal(t, 1, len(points.Points))

				point := points.Points[0]
				assert.Equal(t, "testMeasure", point.Measurement)
				assert.Contains(t, point.Fields, "testCounter.count")
				assert.Equal(t, map[string]string{
					"reporterTag1": "valRep1",
					"reporterTag2": "valRep2",
					"reporterTag3": "valRep3",
					"tag1":         "val1",
					"tag2":         "val2",
				}, point.Tags)

				return nil, nil
			},
		}),
	)
	go reporter.Run()
	defer reporter.Stop()

	metric := NewCounter("testCounter",
		WithMeasurement("testMeasure"),
		WithReporter(reporter),
		WithTags(map[string]string{
			"tag1": "val1",
			"tag2": "val2",
		}),
	)
	for i := 0; i <= 10; i++ {
		time.Sleep(10 * time.Millisecond)

		metric.Inc(int64(i))
	}

	assert.NotEmpty(t, count.Get())
}

func TestNewGauge(t *testing.T) {
	var count = concurrent.NewInt()

	reporter := NewReporter("", "testDB",
		Interval(30*time.Millisecond),
		Tags(map[string]string{
			"reporterTag1": "valRep1",
			"reporterTag2": "valRep2",
			"reporterTag3": "valRep3",
		}),
		Registry(metrics.NewRegistry()),
		withDBClient(&testClient{
			pingCall: func() (duration time.Duration, s string, err error) {
				return 0, "", nil
			},
			writeCall: func(points client.BatchPoints) (response *client.Response, err error) {
				count.Increase()
				assert.Equal(t, "testDB", points.Database)
				assert.True(t, points.Time.Before(time.Now()))
				assert.True(t, points.Time.After(time.Now().Add(-time.Second)))
				assert.Equal(t, 1, len(points.Points))

				point := points.Points[0]
				assert.Equal(t, "testMeasure", point.Measurement)
				assert.Contains(t, point.Fields, "testGauge.gauge")
				assert.Equal(t, map[string]string{
					"reporterTag1": "valRep1",
					"reporterTag2": "valRep2",
					"reporterTag3": "valRep3",
					"tag1":         "val1",
					"tag2":         "val2",
				}, point.Tags)

				return nil, nil
			},
		}),
	)
	go reporter.Run()
	defer reporter.Stop()

	metric := NewGauge("testGauge",
		WithMeasurement("testMeasure"),
		WithReporter(reporter),
		WithTags(map[string]string{
			"tag1": "val1",
			"tag2": "val2",
		}),
	)
	for i := 0; i <= 10; i++ {
		time.Sleep(10 * time.Millisecond)
		metric.Update(int64(i))
	}

	assert.NotEmpty(t, count.Get())
}

func TestNewGaugeFloat64(t *testing.T) {
	var count = concurrent.NewInt()

	reporter := NewReporter("", "testDB",
		Interval(30*time.Millisecond),
		Tags(map[string]string{
			"reporterTag1": "valRep1",
			"reporterTag2": "valRep2",
			"reporterTag3": "valRep3",
		}),
		Registry(metrics.NewRegistry()),
		withDBClient(&testClient{
			pingCall: func() (duration time.Duration, s string, err error) {
				return 0, "", nil
			},
			writeCall: func(points client.BatchPoints) (response *client.Response, err error) {
				count.Increase()
				assert.Equal(t, "testDB", points.Database)
				assert.True(t, points.Time.Before(time.Now()))
				assert.True(t, points.Time.After(time.Now().Add(-time.Second)))
				assert.Equal(t, 1, len(points.Points))

				point := points.Points[0]
				assert.Equal(t, "testMeasure", point.Measurement)
				assert.Contains(t, point.Fields, "testGaugeFloat64.gauge")
				assert.Equal(t, map[string]string{
					"reporterTag1": "valRep1",
					"reporterTag2": "valRep2",
					"reporterTag3": "valRep3",
					"tag1":         "val1",
					"tag2":         "val2",
				}, point.Tags)

				return nil, nil
			},
		}),
	)
	go reporter.Run()
	defer reporter.Stop()

	metric := NewGaugeFloat64("testGaugeFloat64",
		WithMeasurement("testMeasure"),
		WithReporter(reporter),
		WithTags(map[string]string{
			"tag1": "val1",
			"tag2": "val2",
		}),
	)
	for i := 0; i <= 10; i++ {
		time.Sleep(10 * time.Millisecond)
		metric.Update(float64(i))
	}

	assert.NotEmpty(t, count.Get())
}

func TestNewTimer(t *testing.T) {
	var count = concurrent.NewInt()

	reporter := NewReporter("", "testDB",
		Interval(30*time.Millisecond),
		Tags(map[string]string{
			"reporterTag1": "valRep1",
			"reporterTag2": "valRep2",
			"reporterTag3": "valRep3",
		}),
		Registry(metrics.NewRegistry()),
		withDBClient(&testClient{
			pingCall: func() (duration time.Duration, s string, err error) {
				return 0, "", nil
			},
			writeCall: func(points client.BatchPoints) (response *client.Response, err error) {
				count.Increase()
				assert.Equal(t, "testDB", points.Database)
				assert.True(t, points.Time.Before(time.Now()))
				assert.True(t, points.Time.After(time.Now().Add(-time.Second)))
				assert.Equal(t, 16, len(points.Points))

				point := points.Points[0]
				assert.Equal(t, "testMeasure", point.Measurement)
				assert.Contains(t, point.Fields, "testTimer.timer")
				assert.Contains(t, point.Tags, "bucket")
				assert.Equal(t, map[string]string{
					"reporterTag1": "valRep1",
					"reporterTag2": "valRep2",
					"reporterTag3": "valRep3",
					"tag1":         "val1",
					"tag2":         "val2",
				}, bucketLessTags(point.Tags))

				return nil, nil
			},
		}),
	)
	go reporter.Run()
	defer reporter.Stop()

	metric := NewTimer("testTimer",
		WithMeasurement("testMeasure"),
		WithReporter(reporter),
		WithTags(map[string]string{
			"tag1": "val1",
			"tag2": "val2",
		}),
	)
	for i := 0; i <= 10; i++ {
		time.Sleep(10 * time.Millisecond)
		metric.Update(time.Duration(i) * time.Second)
	}

	assert.NotEmpty(t, count.Get())
}

func TestNewMeter(t *testing.T) {
	var count = concurrent.NewInt()

	reporter := NewReporter("", "testDB",
		Interval(30*time.Millisecond),
		Tags(map[string]string{
			"reporterTag1": "valRep1",
			"reporterTag2": "valRep2",
			"reporterTag3": "valRep3",
		}),
		Registry(metrics.NewRegistry()),
		withDBClient(&testClient{
			pingCall: func() (duration time.Duration, s string, err error) {
				return 0, "", nil
			},
			writeCall: func(points client.BatchPoints) (response *client.Response, err error) {
				count.Increase()
				assert.Equal(t, "testDB", points.Database)
				assert.True(t, points.Time.Before(time.Now()))
				assert.True(t, points.Time.After(time.Now().Add(-time.Second)))
				assert.Equal(t, 5, len(points.Points))

				point := points.Points[0]
				assert.Equal(t, "testMeasure", point.Measurement)
				assert.Contains(t, point.Fields, "testMeter.meter")
				assert.Contains(t, point.Tags, "bucket")
				assert.Equal(t, map[string]string{
					"reporterTag1": "valRep1",
					"reporterTag2": "valRep2",
					"reporterTag3": "valRep3",
					"tag1":         "val1",
					"tag2":         "val2",
				}, bucketLessTags(point.Tags))

				return nil, nil
			},
		}),
	)
	go reporter.Run()
	defer reporter.Stop()

	metric := NewMeter("testMeter",
		WithMeasurement("testMeasure"),
		WithReporter(reporter),
		WithTags(map[string]string{
			"tag1": "val1",
			"tag2": "val2",
		}),
	)
	for i := 0; i <= 10; i++ {
		time.Sleep(10 * time.Millisecond)
		metric.Mark(int64(i))
	}

	assert.NotEmpty(t, count.Get())
}

func TestNewHistogram(t *testing.T) {
	var count = concurrent.NewInt()

	reporter := NewReporter("", "testDB",
		Interval(30*time.Millisecond),
		Tags(map[string]string{
			"reporterTag1": "valRep1",
			"reporterTag2": "valRep2",
			"reporterTag3": "valRep3",
		}),
		Registry(metrics.NewRegistry()),
		withDBClient(&testClient{
			pingCall: func() (duration time.Duration, s string, err error) {
				return 0, "", nil
			},
			writeCall: func(points client.BatchPoints) (response *client.Response, err error) {
				count.Increase()
				assert.Equal(t, "testDB", points.Database)
				assert.True(t, points.Time.Before(time.Now()))
				assert.True(t, points.Time.After(time.Now().Add(-time.Second)))
				assert.Equal(t, 12, len(points.Points))

				point := points.Points[0]
				assert.Equal(t, "testMeasure", point.Measurement)
				assert.Contains(t, point.Fields, "testHisto.histogram")
				assert.Contains(t, point.Tags, "bucket")
				assert.Equal(t, map[string]string{
					"reporterTag1": "valRep1",
					"reporterTag2": "valRep2",
					"reporterTag3": "valRep3",
					"tag1":         "val1",
					"tag2":         "val2",
				}, bucketLessTags(point.Tags))

				return nil, nil
			},
		}),
	)
	go reporter.Run()
	defer reporter.Stop()

	metric := NewHistogram("testHisto",
		WithMeasurement("testMeasure"),
		WithReporter(reporter),
		WithTags(map[string]string{
			"tag1": "val1",
			"tag2": "val2",
		}),
	)
	for i := 0; i <= 10; i++ {
		time.Sleep(10 * time.Millisecond)
		metric.Update(int64(i))
	}

	assert.NotEmpty(t, count.Get())
}

func TestNewBasicGauge(t *testing.T) {
	var count = concurrent.NewInt()

	reg := metrics.NewRegistry()
	reporter := NewReporter("", "testDB",
		Interval(30*time.Millisecond),
		Tags(map[string]string{
			"reporterTag1": "valRep1",
			"reporterTag2": "valRep2",
			"reporterTag3": "valRep3",
		}),
		Registry(reg),
		withDBClient(&testClient{
			pingCall: func() (duration time.Duration, s string, err error) {
				return 0, "", nil
			},
			writeCall: func(points client.BatchPoints) (response *client.Response, err error) {
				count.Increase()

				assert.Equal(t, "testDB", points.Database)
				assert.True(t, points.Time.Before(time.Now()))
				assert.True(t, points.Time.After(time.Now().Add(-time.Second)))
				assert.Equal(t, 1, len(points.Points))

				point := points.Points[0]
				assert.Equal(t, "default", point.Measurement)
				assert.Contains(t, point.Fields, "megaMetric.gauge")
				assert.Equal(t, map[string]string(nil), point.Tags)

				return nil, nil
			},
		}),
	)
	go reporter.Run()
	defer reporter.Stop()

	metric := metrics.NewGauge()
	if err := reg.Register("megaMetric", metric); err != nil {
		t.Fatal(err)
	}
	for i := 0; i <= 10; i++ {
		time.Sleep(10 * time.Millisecond)
		metric.Update(int64(i))
	}

	assert.NotEmpty(t, count.Get())
}

type custGauge struct{}

func (s custGauge) Snapshot() metrics.Gauge                 { panic("not implemented") }
func (s custGauge) Update(int64)                            { panic("not implemented") }
func (s custGauge) Value() int64                            { panic("not implemented") }
func (s custGauge) AddPoints([]client.Point) []client.Point { panic("not implemented") }

type custUnknown struct{}

func (s custUnknown) AddPoints([]client.Point) []client.Point { panic("not implemented") }

func TestRegister(t *testing.T) {
	type args struct {
		name    string
		metric  Metric
		options []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "add custom gauge",
			args: args{
				name:    "name1",
				metric:  custGauge{},
				options: nil,
			},
		},
		{
			name: "try adding unknown metric",
			args: args{
				name:    "name2",
				metric:  custUnknown{},
				options: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Register(tt.args.name, tt.args.metric, tt.args.options...)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func bucketLessTags(tags map[string]string) map[string]string {
	m := make(map[string]string, len(tags))
	for k, v := range tags {
		if k == "bucket" {
			continue
		}
		m[k] = v
	}
	return m
}
