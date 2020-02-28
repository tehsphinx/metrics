package metrics

import (
	"context"
	"log"
	"net/url"
	"time"

	client "github.com/influxdata/influxdb1-client"
	"github.com/rcrowley/go-metrics"
)

var defaultReporter Reporter

// SetDefaultReporter sets the default reporter to be used for all metrics. It can be overwritten per Measurement.
func SetDefaultReporter(reporter Reporter) {
	defaultReporter = reporter
}

// Reporter defines a metrics reporter. It is responsible for connection handling
// and sending data to it. It also holds the metrics registry all the metrics get registered to.
// Implementing the Reporter interface is useful for testing or changing the way the reporter behaves.
type Reporter interface {
	Run()
	Register(name string, metric Metric) error
	Get(name string) (Metric, bool)
	Tags() map[string]string
	Stop()
}

type typeChecker func(m metric) bool

// NewReporter creates a new reporter which holds the influxDB connection and sends data to it.
func NewReporter(influxURL, database string, options ...ReporterOption) Reporter {
	dbURL, err := url.Parse(influxURL)
	if err != nil {
		log.Printf("metrics.NewReporter: unable to parse InfluxDB url %s: %v", influxURL, err)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	r := &reporter{
		server: server{
			URL: *dbURL,
			DB:  database,
		},
		registry: metrics.DefaultRegistry,
		interval: 10 * time.Second,
		ctx:      ctx,
		cancel:   cancel,
	}

	for _, option := range options {
		option(r)
	}
	return r
}

type server struct {
	URL  url.URL
	DB   string
	User string
	Pass string
}

// reporter implements a influxDB reporter. This is responsible for the influxDB connection
// and sending data to it. It also holds the metrics registry all the metrics get registered to.
type reporter struct {
	registry metrics.Registry
	client   dbClient
	server   server

	interval time.Duration
	tags     map[string]string
	align    bool

	running bool
	ctx     context.Context
	cancel  context.CancelFunc
}

// Register registers a metric to the reporter. Data points from a registered
// metric will be collected via the AddPoints endpoint and then sent to influxDB.
//
// This function is only needed for custom metrics. Functions creating a metric
// in this package register themselves.
func (r *reporter) Register(name string, metric Metric) error {
	return r.registry.Register(name, metric)
}

// Get returns a metric by name. Returns false if it does not exist or does not implement Metric.
func (r *reporter) Get(name string) (Metric, bool) {
	m, ok := r.registry.Get(name).(Metric)
	return m, ok
}

// Tags returns the tags. The return value should not be modified.
func (r *reporter) Tags() map[string]string {
	return r.tags
}

// Run starts sending measurements regularly with given interval.
// This is a blocking call and is usually called with `go reporter.Run()`.
func (r *reporter) Run() {
	if r.running {
		log.Println("metrics.Reporter already running")
		return
	}

	if err := r.open(); err != nil {
		log.Printf("unable to reconnect InfluxDB client: %v", err)
		return
	}

	r.run()
}
func (r *reporter) run() {
	var (
		pts            []client.Point
		intervalTicker = time.NewTicker(r.interval)
		pingTicker     = time.NewTicker(time.Second * 5)
	)

	for {
		select {
		case <-r.ctx.Done():
			return
		case <-intervalTicker.C:
			pts = pts[:0]
			pts = r.getPoints(pts)

			if err := r.write(pts); err != nil {
				log.Printf("unable to send metrics to InfluxDB: %v", err)
			}
		case <-pingTicker.C:
			_, _, err := r.client.Ping()
			if err != nil {
				log.Printf("got error while sending a ping to InfluxDB: %v", err)

				if err = r.open(); err != nil {
					log.Printf("unable to reconnect InfluxDB client: %v", err)
				}
			}
		}
	}
}

func (r *reporter) getPoints(pts []client.Point) []client.Point {
	r.registry.Each(func(name string, data interface{}) {
		m, ok := data.(Metric)
		if !ok {
			pts = r.basicMetric(pts, name, data)
			return
		}

		pts = m.AddPoints(pts)
	})
	return pts
}

func (r *reporter) basicMetric(pts []client.Point, name string, data interface{}) []client.Point {
	var m Metric
	switch metric := data.(type) {
	case metrics.Counter:
		m = newCounter(name, WithMetric(metric))
	case metrics.Gauge:
		m = newGauge(name, WithMetric(metric))
	case metrics.GaugeFloat64:
		m = newGaugeFloat64(name, WithMetric(metric))
	case metrics.Histogram:
		m = newHistogram(name, WithMetric(metric))
	case metrics.Meter:
		m = newMeter(name, WithMetric(metric))
	case metrics.Timer:
		m = newTimer(name, WithMetric(metric))
	default:
		return pts
	}

	return m.AddPoints(pts)
}

func (r *reporter) open() (err error) {
	if r.client != nil {
		return nil
	}

	r.client, err = client.NewClient(client.Config{
		URL:      r.server.URL,
		Username: r.server.User,
		Password: r.server.Pass,
	})

	return err
}

func (r *reporter) write(points []client.Point) error {
	bps := client.BatchPoints{
		Points:   points,
		Database: r.server.DB,
		Time:     r.getNow(),
	}

	_, err := r.client.Write(bps)
	return err
}

func (r *reporter) getNow() time.Time {
	now := time.Now()
	if r.align {
		now = now.Truncate(r.interval)
	}
	return now
}

// Stop stops the reporter. It should be discarded after and cannot be restartet.
func (r *reporter) Stop() {
	r.cancel()
}
