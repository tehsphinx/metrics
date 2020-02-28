package metrics

import (
	"testing"
	"time"

	client "github.com/influxdata/influxdb1-client"
	"github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"
)

func Test_timer_AddPoints(t *testing.T) {
	type fields struct {
		name    string
		measure string
		tags    map[string]string
	}
	type args struct {
		pts []client.Point
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantLen int
	}{
		{
			name: "test AddPoints",
			fields: fields{
				measure: "measure1",
				name:    "testName",
				tags:    map[string]string{"fooby": "bar"},
			},
			args:    args{},
			wantLen: 16,
		},
	}
	suffix := suffTimer

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric := NewTimer(tt.fields.name, WithTags(tt.fields.tags), WithMeasurement(tt.fields.measure))
			metric.Update(5 * time.Second)

			s := metrics.DefaultRegistry.Get(tt.fields.measure + "/" + tt.fields.name + suffix).(Metric)
			got := s.AddPoints(tt.args.pts)

			assert.Equal(t, tt.wantLen, len(got))
			for _, point := range got {
				for k, v := range tt.fields.tags {
					assert.Contains(t, point.Tags, k)
					assert.Contains(t, point.Tags[k], v)
				}
				assert.Contains(t, point.Tags, "bucket")
				assert.Equal(t, len(tt.fields.tags)+1, len(point.Tags))
				assert.Equal(t, tt.fields.measure, point.Measurement)
				assert.Contains(t, point.Fields, tt.fields.name+suffix)
			}
		})
	}
}
