package metrics

import (
	client "github.com/influxdata/influxdb1-client"
)

func buildBucketVals(buckets []string, field string) map[string]map[string]interface{} {
	var m = make(map[string]map[string]interface{}, len(buckets))
	for _, bucket := range buckets {
		m[bucket] = map[string]interface{}{
			field: 0.0,
		}
	}
	return m
}

func buildBucketTags(buckets []string, tags map[string]string) map[string]map[string]string {
	var m = make(map[string]map[string]string, len(buckets))
	for _, bucket := range buckets {
		m[bucket] = bucketTags(bucket, tags)
	}
	return m
}

func bucketTags(bucket string, tags map[string]string) map[string]string {
	if bucket == "" {
		return tags
	}

	m := make(map[string]string, len(tags)+1)
	for tk, tv := range tags {
		m[tk] = tv
	}
	m["bucket"] = bucket
	return m
}

func composeTags(reporterTags, metricTags map[string]string) map[string]string {
	m := make(map[string]string, len(reporterTags)+len(metricTags))
	for k, v := range reporterTags {
		m[k] = v
	}
	for k, v := range metricTags {
		m[k] = v
	}
	return m
}

func getPoint(measurement string, fields map[string]interface{}, tags map[string]string) client.Point {
	return client.Point{
		Measurement: measurement,
		Tags:        tags,
		Fields:      fields,
	}
}
