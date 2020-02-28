package metrics

import (
	"log"
	"reflect"
	"strconv"
	"sync"

	client "github.com/influxdata/influxdb1-client"
	"github.com/rcrowley/go-metrics"
)

const (
	suffGauge     = ".gauge"
	suffTimer     = ".timer"
	suffCounter   = ".count"
	suffHistogram = ".histogram"
	suffMeter     = ".meter"
)

// Metric defines an interface to be implemented for custom metrics.
//
// For a metric to be used, it needs to be registered.
// If no registry was provided to the Reporter, it uses the default registry
// and can be registered with `metrics.Register(name, metric)`.
// If a custom registry was provided to the Reporter, either register it to
// that registry or use the Reporter.Register function to register it.
type Metric interface {
	AddPoints(pts []client.Point) []client.Point
}

type metric interface {
	regName() string
	regIncr()
	AddPoints(pts []client.Point) []client.Point
}

func newMetric(name string, options ...Option) *baseMetric {
	m := &baseMetric{
		name:        name,
		reporter:    defaultReporter,
		measurement: "default",
		suffix:      ".metric",
		regMutex:    &sync.Mutex{},
	}
	for _, option := range options {
		option(m)
	}
	if m.reporter != nil {
		m.tags = composeTags(m.reporter.Tags(), m.tags)
	}
	return m
}

type baseMetric struct {
	name string
	incr int

	metric interface{}

	reporter    Reporter
	measurement string
	tags        map[string]string
	suffix      string

	regMutex *sync.Mutex
}

func (s *baseMetric) regName() string {
	suffix := s.suffix
	if s.incr != 0 {
		suffix = strconv.Itoa(s.incr) + s.suffix
	}
	return s.measurement + "/" + s.name + suffix
}
func (s *baseMetric) regIncr() {
	s.incr++
}

func (s *baseMetric) register(m metric, typeCheck typeChecker) Metric {
	s.regMutex.Lock()
	defer s.regMutex.Unlock()

	return s.reg(m, typeCheck)
}
func (s *baseMetric) reg(m metric, typeCheck typeChecker) Metric {
	regName := m.regName()
	if s.reporter == nil {
		log.Println("WARNING: no (default) metrics reporter set")
		if err := metrics.Register(regName, m); err != nil {
			log.Printf("default registry: metric could not be registered: %v", err)
		}
		return m
	}

	if mtrx, ok := s.reporter.Get(regName); ok {
		if reflect.TypeOf(m) == reflect.TypeOf(mtrx) {
			return mtrx
		}
		m.regIncr()
		return s.reg(m, typeCheck)
	}

	if err := s.reporter.Register(regName, m); err != nil {
		// this should never happen because of the lock above
		log.Printf("metric could not be registered: %v", err)
	}
	return m
}

const (
	count    = "count"
	max      = "max"
	mean     = "mean"
	min      = "min"
	p50      = "p50"
	p75      = "p75"
	p95      = "p95"
	p99      = "p99"
	p999     = "p999"
	p9999    = "p9999"
	stddev   = "stddev"
	variance = "variance"
	m1       = "m1"
	m5       = "m5"
	m15      = "m15"
	meanrate = "meanrate"
)
