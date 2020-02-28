package metrics

import (
	"net/url"
	"testing"
	"time"

	client "github.com/influxdata/influxdb1-client"
	"github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"
)

func TestSetDefaultReporter(t *testing.T) {
	reporter := NewReporter("", "")
	SetDefaultReporter(reporter)

	assert.Equal(t, reporter, defaultReporter)
}

func TestNewReporter(t *testing.T) {
	type args struct {
		influxURL string
		database  string
		options   []ReporterOption
	}
	tests := []struct {
		name string
		args args
		want Reporter
	}{
		{
			name: "basic",
			args: args{},
			want: &reporter{
				registry: metrics.DefaultRegistry,
				interval: 10 * time.Second,
			},
		},
		{
			name: "basic with data",
			args: args{
				influxURL: "https://someDomain.com/somePath",
				database:  "someDB",
			},
			want: &reporter{
				registry: metrics.DefaultRegistry,
				interval: 10 * time.Second,
				server: server{
					URL: url.URL{
						Scheme: "https",
						Host:   "someDomain.com",
						Path:   "/somePath",
					},
					DB: "someDB",
				},
			},
		},
		{
			name: "tags",
			args: args{
				options: []ReporterOption{
					Tags(map[string]string{"foo": "bar"}),
				},
			},
			want: &reporter{
				registry: metrics.DefaultRegistry,
				interval: 10 * time.Second,
				tags:     map[string]string{"foo": "bar"},
			},
		},
		{
			name: "registry",
			args: args{
				options: []ReporterOption{
					Registry(metrics.NewRegistry()),
				},
			},
			want: &reporter{
				registry: metrics.NewRegistry(),
				interval: 10 * time.Second,
			},
		},
		{
			name: "interval",
			args: args{
				options: []ReporterOption{
					Interval(5 * time.Second),
				},
			},
			want: &reporter{
				registry: metrics.DefaultRegistry,
				interval: 5 * time.Second,
			},
		},
		{
			name: "align",
			args: args{
				options: []ReporterOption{
					Align(),
				},
			},
			want: &reporter{
				registry: metrics.DefaultRegistry,
				interval: 10 * time.Second,
				align:    true,
			},
		},
		{
			name: "auth",
			args: args{
				options: []ReporterOption{
					Auth("user", "pass"),
				},
			},
			want: &reporter{
				registry: metrics.DefaultRegistry,
				interval: 10 * time.Second,
				server: server{
					User: "user",
					Pass: "pass",
				},
			},
		},
		{
			name: "clientDB",
			args: args{
				options: []ReporterOption{
					withDBClient(&testClient{}),
				},
			},
			want: &reporter{
				registry: metrics.DefaultRegistry,
				interval: 10 * time.Second,
				client:   &testClient{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewReporter(tt.args.influxURL, tt.args.database, tt.args.options...).(*reporter)
			got.ctx = nil
			got.cancel = nil
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_reporter_Get(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name     string
		existing map[string]interface{}
		args     args
		want     Metric
		wantOk   bool
	}{
		{
			name: "get non existing item",
			existing: map[string]interface{}{
				"foo": &timer{},
				"bar": 8,
			},
			args: args{
				name: "non-existing",
			},
			want:   nil,
			wantOk: false,
		},
		{
			name: "try getting a non-Metric item",
			existing: map[string]interface{}{
				"foo": &timer{},
				"bar": 8,
			},
			args: args{
				name: "bar",
			},
			want:   nil,
			wantOk: false,
		},
		{
			name: "get a metric",
			existing: map[string]interface{}{
				"foo": &timer{},
				"bar": 8,
			},
			args: args{
				name: "foo",
			},
			want:   &timer{},
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := metrics.NewRegistry()
			r := NewReporter("", "testDB", Registry(reg))
			for name, item := range tt.existing {
				if err := reg.Register(name, item); err != nil {
					t.Error(err)
				}
			}

			got, ok := r.Get(tt.args.name)
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Benchmark_reporter_getPoints(b *testing.B) {
	r := NewReporter("", "", Registry(metrics.NewRegistry())).(*reporter)

	NewTimer("metric", WithReporter(r)).Update(5)
	NewGauge("metric", WithReporter(r)).Update(5)
	NewGaugeFloat64("metric", WithReporter(r)).Update(5.542)
	NewCounter("metric", WithReporter(r)).Inc(5)
	NewHistogram("metric", WithReporter(r)).Update(5)
	NewMeter("metric", WithReporter(r)).Mark(5)

	b.ResetTimer()

	var pts []client.Point
	for i := 0; i < b.N; i++ {
		pts = pts[:0]

		pts = r.getPoints(pts)
	}
}
