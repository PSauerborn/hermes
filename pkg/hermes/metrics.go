package hermes

import (
    "fmt"
    "errors"
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    log "github.com/sirupsen/logrus"

    "github.com/PSauerborn/hermes/pkg/utils"
)

var (
    Config *HermesConfig

    // define maps used to store metrics
    Gauges     = map[string]*prometheus.GaugeVec{}
    Counters   = map[string]*prometheus.CounterVec{}
    Histograms = map[string]*prometheus.HistogramVec{}
    Summaries  = map[string]*prometheus.SummaryVec{}

    // define custom errors for application
    ErrInvalidGauge          = errors.New("Invalid gauge configuration")
    ErrInvalidCounter        = errors.New("Invalid gauge configuration")
    ErrUnregisteredMetric    = errors.New("Unregistered metric")
    ErrInvalidGaugeOperation = errors.New("Invalid gauge operation")
    ErrInvalidLabels         = errors.New("Invalid label configuration")
)

// function used to start new prometheus server
// to scrape metrics from Hermes
func ListenPrometheus(config HermesConfig, listenPort int) {
    // create prometheus metric objects from configuration
    InitializeMetrics(config)
    // create http interface to listen for prometheus scrape jobs
    connectionString := fmt.Sprintf(":%d", listenPort)
    http.Handle("/metrics", promhttp.Handler())
    log.Fatal(http.ListenAndServe(connectionString, nil))
}

// function used to initialize hermes metrics by iterating
// over the JSON configuration file and generating prometheus
// Gauges/Counters for all the specified metrics
func InitializeMetrics(config HermesConfig) error {
    Config = &config
    // create gauges from config
    for _, gauge := range(config.Gauges) {
        log.Debug(fmt.Sprintf("creating new gauge from config %+v", gauge))
        if err := NewGauge(gauge); err != nil {
            log.Fatal(fmt.Errorf("unable to create new gauge: %v", err))
        }
    }
    // create counters from config
    for _, counter := range(config.Counters) {
        log.Debug(fmt.Sprintf("creating new counter from config %+v", counter))
        if err := NewCounter(counter); err != nil {
            log.Fatal(fmt.Errorf("unable to create new counter: %v", err))
        }
    }
    // create histograms from config
    for _, histogram := range(config.Histograms) {
        log.Debug(fmt.Sprintf("creating new histogram from config %+v", histogram))
        if err := NewHistogram(histogram); err != nil {
            log.Fatal(fmt.Errorf("unable to create new histogram: %v", err))
        }
    }
    // create summaries from config
    for _, summary := range(config.Summaries) {
        log.Debug(fmt.Sprintf("creating new summary from config %+v", summary))
        if err := NewSummary(summary); err != nil {
            log.Fatal(fmt.Errorf("unable to create new summary: %v", err))
        }
    }
    return nil
}

// function used to determine the metric type
// based on a particular metric name
func GetMetricType(metric string) (string, error) {
    // check if metric is present in registered counters
    if _, ok := Counters[metric]; ok {
        return "counter", nil
    }
    // check if metric is present in registered gauges
    if _, ok := Gauges[metric]; ok {
        return "gauge", nil
    }
    // check if metric is present in registered histograms
    if _, ok := Histograms[metric]; ok {
        return "histogram", nil
    }
    // check if metric is present in registered summaries
    if _, ok := Summaries[metric]; ok {
        return "summary", nil
    }
    return "", ErrUnregisteredMetric
}

// function used to determine if a given set of labels
// matches the label configuration expected for the
// specified metric
func IsValidLabelConfig(receivedLabels map[string]string, expectedLabels []string) bool {
    // iterate over keys of given labels and check if key
    // is present in labels defined in JSON config
    for key := range(receivedLabels) {
        // return false if any supplied label is NOT in defined config
        if !utils.SliceContains(expectedLabels, key) {
            return false
        }
    }
    return true
}

// function used to convert labels into prometheus.Labels instance
// by filtering out the lables that are included both on the global
// hermes configuration file and the JSON from the UDP packet
func SetPrometheusLabels(labels map[string]string, labelConfig []string) (prometheus.Labels, error) {

    // check that labels from payload match labels defined in config
    if !IsValidLabelConfig(labels, labelConfig) {
        log.Error(fmt.Sprintf("invalid label configuration. expecting %d but received %d",
            len(labels), len(labelConfig)))
        return prometheus.Labels{}, ErrInvalidLabels
    }
    // generate prometheus label instance and add labels given in payload
    promLabels := prometheus.Labels{}
    for _, metricLabel := range(labelConfig) {
        if label, ok := labels[metricLabel]; ok {
            promLabels[metricLabel] = label
        }
    }
    return promLabels, nil
}

// function used to generate prometheus labels based on config.
// note that the labels provided in the UDP packet are not set
// on the counter/gauge unless they have also been defined in
// the JSON config file
func GenerateLabels(labels map[string]string, metricType,
    metricName string) (prometheus.Labels, error) {

    var (promLabels prometheus.Labels; err error)
    // retrieve labels registered for metric in config
    // and set against lables provided in payload
    switch metricType {
    case "counter":
        for _, counter := range(Config.Counters) {
            if counter.MetricName == metricName {
                // create labels for counter instance
                promLabels, err = SetPrometheusLabels(labels, counter.Labels)
            }
        }
    case "gauge":
        for _, gauge := range(Config.Gauges) {
            if gauge.MetricName == metricName {
                // create labels for gauge instance
                promLabels, err = SetPrometheusLabels(labels, gauge.Labels)
            }
        }
    case "histogram":
        for _, histogram := range(Config.Histograms) {
            if histogram.MetricName == metricName {
                // create labels for histogram instance
                promLabels, err = SetPrometheusLabels(labels, histogram.Labels)
            }
        }
    case "summary":
        for _, summary := range(Config.Summaries) {
            if summary.MetricName == metricName {
                // create labels for summary instance
                promLabels, err = SetPrometheusLabels(labels, summary.Labels)
            }
        }
    }
    return promLabels, err
}
