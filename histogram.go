package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// function used to make an observation on a particular histogram
func ObserveHistogram(name string, histogramJson HistogramJSON) error {
	if histogram, ok := Histograms[name]; ok {
        log.Info(fmt.Sprintf("making histogram observation %f on '%s' %v", histogramJson.Observation, name, histogram))
        // generate labels for prometheus metric and check for errors
        labels, err := GenerateLabels(histogramJson.Labels, "histogram", name)
        if err != nil {
            return err
        }
        histogram.With(labels).Observe(histogramJson.Observation)
        return nil
    }
    return ErrUnregisteredMetric
}

// function used to create a new histogram instance. Pointers to the
// prometheus histograms are stored in the global histogram map, which
// maps the name of the histogram/metric to the prometheus pointer
// that stores the metrics themselves
func NewHistogram(histogram HermesHistogram) error {
	opts := prometheus.HistogramOpts{Name: histogram.MetricName, Help: histogram.MetricDescription}
	// create new histogram instance
	promHistogram := prometheus.NewHistogramVec(opts, histogram.Labels)
	// register gauge and insert into maps
    prometheus.MustRegister(promHistogram)
    Histograms[histogram.MetricName] = promHistogram
    return nil
}