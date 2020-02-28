package metrics

import (
	"time"

	client "github.com/influxdata/influxdb1-client"
)

type dbClient interface {
	Write(points client.BatchPoints) (*client.Response, error)
	Ping() (time.Duration, string, error)
}
