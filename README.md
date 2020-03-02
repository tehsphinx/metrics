# InfluxDB Metrics Reporter

Package metrics is based on [github.com/rcrowley/go-metrics](https://github.com/rcrowley/go-metrics) and loosely inspired
by [github.com/vrischmann/go-metrics-influxdb](https://github.com/vrischmann/go-metrics-influxdb).


## Why?

Why write a new package if there already is a influxDB reporter?

The main difference between this package and `vrischmann/go-metrics-influxdb` is the API.
This API is based on the `go-metrics` api and separates the registry (built into the Reporter) from the measurement/metrics.
This enables a global Reporter and the creation of independant metrics distributed across
the codebase without passing the Reporter around. A reporter can however be injected into a metrics if needed.

With `vrischmann/go-metrics-influxdb` a reporter can only handle one measurement, 
meaning a new reporter is needed per measurement. This entails more goroutines and 
more calls to the influxDB in larger applications with multiple measurements.

## Usage

### Global Reporter, metrics accross the application

```go
func main() {
	// Initialize reporter and set it as default (global).
	rep := metrics.NewReporter("http://localhost:8086", "metrics", metrics.Interval(1*time.Second))
	metrics.SetDefaultReporter(rep)

	// Start reporting all registered or to be registered metrics.
	go rep.Run()
}
```

This creates a new reporter, sets it as the default reporter and starts it.
The reporter will then start a Timer with given interval (default: 10 seconds).
It will collect and report data from all metrics registered to the default reporter.

To register a metric use the following code anywhere in your server:

```go
func registerGauge() {
	// Create and register a new gauge metric.
	metric := metrics.NewGauge("MetricName", metrics.WithMeasurement("MeasurementName"))
}
```

To update the data of the metric:

```go
// Update the gauge metric value.
metric.Update(n)
```

For working code see the [simple-gauge example](examples/simple_gauge/main.go).

### Without Global Reporter

Some might prefer to inject the reporter with the metric instead of using a global reporter.
This is supported with the WithReporter option when creating a metric:

```go
rep := metrics.NewReporter("http://localhost:8086", "metrics")

// inject the reporter to be registered to:
metric := metrics.NewGauge("MetricName", metrics.WithMeasurement("MeasurementName"), metrics.WithReporter(rep))
```

For working code see the [multiple-reporters example](examples/multiple_reporters/main.go).

:warning: When creating multiple repoters use the `Registry` option to supply a new registry to each reporter.
If multiple reporters are created without this option, they all use the default registry. 
Meaning all reporters will report data from all the metrics.

### Metrics

New metrics can be created with the `metrics.NewXY` functions.

For more information on the different metric types see [go-metrics](https://github.com/rcrowley/go-metrics).
If a `go-metrics` metric is not implemented here, please open an issue.
