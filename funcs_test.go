package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildBucketTags(t *testing.T) {
	type args struct {
		buckets []string
		tags    map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]map[string]string
	}{
		{
			name: "",
			args: args{
				buckets: []string{"buck1", "buck2"},
				tags:    map[string]string{"tag1": "val1", "tag2": "val2"},
			},
			want: map[string]map[string]string{
				"buck1": {"tag1": "val1", "tag2": "val2", "bucket": "buck1"},
				"buck2": {"tag1": "val1", "tag2": "val2", "bucket": "buck2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildBucketTags(tt.args.buckets, tt.args.tags)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_bucketTags(t *testing.T) {
	type args struct {
		bucket string
		tags   map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "add bucket to tags",
			args: args{
				bucket: "bucket",
				tags:   map[string]string{"foo1": "bar1", "foo3": "bar3"},
			},
			want: map[string]string{"foo1": "bar1", "foo3": "bar3", "bucket": "bucket"},
		},
		{
			name: "add empty bucket to tags",
			args: args{
				bucket: "",
				tags:   map[string]string{"foo1": "bar1", "foo3": "bar3"},
			},
			want: map[string]string{"foo1": "bar1", "foo3": "bar3"},
		},
		{
			name: "add bucket to empty tags list",
			args: args{
				bucket: "buckFoo",
				tags:   nil,
			},
			want: map[string]string{"bucket": "buckFoo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bucketTags(tt.args.bucket, tt.args.tags)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_composeTags(t *testing.T) {
	type args struct {
		reporterTags map[string]string
		metricTags   map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "compose tag lists",
			args: args{
				reporterTags: map[string]string{"foo": "bar", "foo1": "bar", "foo2": "bar2"},
				metricTags:   map[string]string{"foo1": "bar1", "foo3": "bar3"},
			},
			want: map[string]string{"foo": "bar", "foo1": "bar1", "foo2": "bar2", "foo3": "bar3"},
		},
		{
			name: "reporter list is nil",
			args: args{
				reporterTags: nil,
				metricTags:   map[string]string{"foo1": "bar1", "foo3": "bar3"},
			},
			want: map[string]string{"foo1": "bar1", "foo3": "bar3"},
		},
		{
			name: "metric list is nil",
			args: args{
				reporterTags: map[string]string{"foo1": "bar1", "foo3": "bar3"},
				metricTags:   nil,
			},
			want: map[string]string{"foo1": "bar1", "foo3": "bar3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := composeTags(tt.args.reporterTags, tt.args.metricTags)
			assert.Equal(t, tt.want, got)
		})
	}
}
