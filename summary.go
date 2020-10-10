package main

import (
    "fmt"
    "github.com/prometheus/client_golang/prometheus"
    log "github.com/sirupsen/logrus"
)

// function used to make an observation on a particular histogram
func ObserveSummary(name string, summaryJson SummaryJSON) error {
    if summary, ok := Summaries[name]; ok {
        log.Info(fmt.Sprintf("making histogram observation %f on '%s' %v", summaryJson.Observation, name, summary))
        // generate labels for prometheus metric and check for errors
        labels, err := GenerateLabels(summaryJson.Labels, "summary", name)
        if err != nil {
            return err
        }
        summary.With(labels).Observe(summaryJson.Observation)
        return nil
    }
    return ErrUnregisteredMetric
}

// function used to create a new histogram instance. Pointers to the
// prometheus histograms are stored in the global histogram map, which
// maps the name of the histogram/metric to the prometheus pointer
// that stores the metrics themselves
func NewSummary(summary HermesSummary) error {
    opts := prometheus.SummaryOpts{Name: summary.MetricName, Help: summary.MetricDescription}
    // create new histogram instance
    promSummary := prometheus.NewSummaryVec(opts, summary.Labels)
    // register gauge and insert into maps
    prometheus.MustRegister(promSummary)
    Summaries[summary.MetricName] = promSummary
    return nil
}