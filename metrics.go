package main

import (
    "fmt"
    "errors"
    "net/http"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    log "github.com/sirupsen/logrus"
)

var (
    Gauges = map[string]*prometheus.GaugeVec{}
    Counters = map[string]*prometheus.CounterVec{}
    Config *HermesConfig
    ErrInvalidGauge = errors.New("invalid gauge configuration")
    ErrInvalidCounter = errors.New("invalid gauge configuration")
    ErrUnregisteredMetric = errors.New("unregistered metric")
    ErrInvalidGaugeOperation = errors.New("invalid gauge operation")
)

// function used to start new prometheus server
// to scrape metrics from Hermes
func ListenPrometheus(config HermesConfig) {
    // create prometheus metric objects from configuration
    InitializeMetrics(config)

    http.Handle("/metrics", promhttp.Handler())
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// function used to determine the metric type i.e.
// based on a particular metric name
func GetMetricType(metric string) *string {
    // check if metric is present in registered counters
    if _, ok := Counters[metric]; ok {
        metricType := "counter"
        return &metricType
    }
    // check if metric is present in registered gauges
    if _, ok := Gauges[metric]; ok {
        metricType := "gauge"
        return &metricType
    }
    return nil
}

// function used to convert labels into prometheus.Labels instance
// by filtering out the lables that are included both on the global
// hermes configuration file and the JSON from the UDP packet
func SetPrometheusLabels(labels map[string]string, labelConfig []string) prometheus.Labels {
    promLabels := prometheus.Labels{}
    for _, metricLabel := range(labelConfig) {
        if label, ok := labels[metricLabel]; ok {
            promLabels[metricLabel] = label
        }
    }
    return promLabels
}

// function used to generate prometheus labels based on config.
// note that the labels provided in the UDP packet are not set
// on the counter/gauge unless they have also been defined in
// the JSON config file
func GenerateLabels(labels map[string]string, metricType, metricName string) prometheus.Labels {
    var promLabels prometheus.Labels
    switch metricType {
    case "counter":
        for _, counter := range(Config.Counters) {
            if counter.MetricName == metricName {
                promLabels = SetPrometheusLabels(labels, counter.Labels)
            }
        }
    case "gauge":
        for _, gauge := range(Config.Gauges) {
            if gauge.MetricName == metricName {
                // create labels for counter instance
                promLabels = SetPrometheusLabels(labels, gauge.Labels)
            }
        }
    }
    return promLabels
}

// function used to increment a particular counter
func IncrementCounter(name string, counterJson CounterJSON) error {
    if counter, ok := Counters[name]; ok {
        log.Info(fmt.Sprintf("incrementing counter '%s' %v", name, counter))
        labels := GenerateLabels(counterJson.Labels, "counter", name)
        counter.With(labels).Inc()
        return nil
    }
    return ErrUnregisteredMetric
}

// function used to set the value on a particular gauge
func SetGauge(name string, gaugeJson GaugeJSON) error {
    if gauge, ok := Gauges[name]; ok {
        log.Info(fmt.Sprintf("setting gauge '%s' %v", name, gauge))
        labels := GenerateLabels(gaugeJson.Labels, "gauge", name)
        gauge.With(labels).Set(*gaugeJson.Value)
        return nil
    }
    return ErrUnregisteredMetric
}

// function used to increment gauge a particular gauge value
func IncrementGauge(name string, gaugeJson GaugeJSON) error {
    if gauge, ok := Gauges[name]; ok {
        log.Info(fmt.Sprintf("incrementing gauge '%s' %v", name, gauge))
        labels := GenerateLabels(gaugeJson.Labels, "gauge", name)
        gauge.With(labels).Inc()
        return nil
    }
    return ErrUnregisteredMetric
}

// function used to decrement a particular gauge value
func DecrementGauge(name string, gaugeJson GaugeJSON) error {
    if gauge, ok := Gauges[name]; ok {
        log.Info(fmt.Sprintf("decrementing gauge '%s' %v", name, gauge))
        labels := GenerateLabels(gaugeJson.Labels, "gauge", name)
        gauge.With(labels).Dec()
        return nil
    }
    return ErrUnregisteredMetric
}

// function used to call correct handler for gauge operations.
// currently, gauge operations support incrementing, deprecating,
// and setting of values.
func ProcessGauge(name string, gaugeJson GaugeJSON) error {
    switch gaugeJson.Operation {
        // increment gauge
    case "increment":
        return IncrementGauge(name, gaugeJson)
        // decrement gauge
    case "decrement":
        return DecrementGauge(name, gaugeJson)
    case "set":
        // ensure that values has been specified if setting gauge
        if gaugeJson.Value != nil {
            return SetGauge(name, gaugeJson)
        } else {
            log.Error(fmt.Sprintf("gauge cannot be set without value"))
            return ErrInvalidGaugeOperation
        }
    default:
        log.Error(fmt.Sprintf("received invalid gauge operation '%s'", gaugeJson.Operation))
        return ErrInvalidGaugeOperation
    }
}

// function used to create new gauge instance. Pointers to the
// prometheus gauges are stored in the global Gauges map, which
// maps the name of the gauge/metric to the prometheus pointer
// that stores the metrics themselves
func NewGauge(gauge HermesGauge) error {
    opts := prometheus.GaugeOpts{Name: gauge.MetricName, Help: gauge.MetricDescription}
    // create new prometheus gauge
    promGauge := prometheus.NewGaugeVec(opts, gauge.Labels)
    // register gauge and insert into maps
    prometheus.MustRegister(promGauge)
    Gauges[gauge.MetricName] = promGauge
    return nil
}

// function used to create new counter instance. Pointers to the
// prometheus counters are stored in the global Gauges map, which
// maps the name of the counter/metric to the prometheus pointer
// that stores the metrics themselves
func NewCounter(counter HermesCounter) error {
    opts := prometheus.CounterOpts{Name: counter.MetricName, Help: counter.MetricDescription}
    // create new counter
    promCounter := prometheus.NewCounterVec(opts, counter.Labels)
    // register counter and insert into maps
    prometheus.MustRegister(promCounter)
    Counters[counter.MetricName] = promCounter
    return nil
}

// function used to initialize hermes metrics by iterating
// over the JSON configuration file and generating prometheus
// Gauges/Counters for all the specified metrics
func InitializeMetrics(config HermesConfig) error {
    Config = &config
    // create gauges from config
    for _, gauge := range(config.Gauges) {
        log.Debug(fmt.Sprintf("creating new gauge from config %+v", gauge))
        err := NewGauge(gauge)
        if err != nil {
            log.Fatal(fmt.Errorf("unable to create new gauge: %v", err))
        }
    }
    // create counters from config
    for _, counter := range(config.Counters) {
        log.Debug(fmt.Sprintf("creating new counter from config %+v", counter))
        err := NewCounter(counter)
        if err != nil {
            log.Fatal(fmt.Errorf("unable to create new counter: %v", err))
        }
    }
    return nil
}

