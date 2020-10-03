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
)

func ListenPrometheus(config HermesConfig) {
    // create prometheus metric objects from configuration
    InitializeMetrics(config)

    http.Handle("/metrics", promhttp.Handler())
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// function used to determin metric type
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
func SetPrometheusLabels(labels map[string]string, labelConfig []string) prometheus.Labels {
    promLabels := prometheus.Labels{}
    for _, metricLabel := range(labelConfig) {
        if label, ok := labels[metricLabel]; ok {
            promLabels[metricLabel] = label
        }
    }
    return promLabels
}

// function used to generate prometheus labels based on config
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
                // create counter labels
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

// function used to set a particular gauge
func SetGauge(name string, gaugeJson GaugeJSON) error {
    if gauge, ok := Gauges[name]; ok {
        log.Info(fmt.Sprintf("setting gauge '%s' %v", name, gauge))
        labels := GenerateLabels(gaugeJson.Labels, "gauge", name)
        gauge.With(labels).Set(gaugeJson.Value)
        return nil
    }
    return ErrUnregisteredMetric
}

// function used to create new gauge
func NewGauge(gauge HermesGauge) error {
    opts := prometheus.GaugeOpts{Name: gauge.MetricName, Help: gauge.MetricDescription}
    // create new prometheus gauge
    promGauge := prometheus.NewGaugeVec(opts, gauge.Labels)
    // register gauge and insert into maps
    prometheus.MustRegister(promGauge)
    Gauges[gauge.MetricName] = promGauge
    return nil
}

// function use to create new counter
func NewCounter(counter HermesCounter) error {
    opts := prometheus.CounterOpts{Name: counter.MetricName, Help: counter.MetricDescription}
    // create new counter
    promCounter := prometheus.NewCounterVec(opts, counter.Labels)
    // register counter and insert into maps
    prometheus.MustRegister(promCounter)
    Counters[counter.MetricName] = promCounter
    return nil
}

// function used to initialize hermes metrics
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

