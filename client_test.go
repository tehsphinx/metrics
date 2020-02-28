package metrics

import (
	"time"

	client "github.com/influxdata/influxdb1-client"
)

type testClient struct {
	writeCall func(points client.BatchPoints) (*client.Response, error)
	pingCall  func() (time.Duration, string, error)
}

func (s *testClient) Write(points client.BatchPoints) (*client.Response, error) {
	return s.writeCall(points)
}

func (s *testClient) Ping() (time.Duration, string, error) {
	return s.pingCall()
}
