package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// function used to set the value on a particular gauge
func SetGauge(name string, gaugeJson GaugeJSON) error {
    if gauge, ok := Gauges[name]; ok {
        log.Info(fmt.Sprintf("setting gauge '%s' %v", name, gauge))
        // generate labels for prometheus metric and check for errors
        labels, err := GenerateLabels(gaugeJson.Labels, "gauge", name)
        if err != nil {
            return err
        }
        gauge.With(labels).Set(*gaugeJson.Value)
        return nil
    }
    return ErrUnregisteredMetric
}

// function used to increment gauge a particular gauge value
func IncrementGauge(name string, gaugeJson GaugeJSON) error {
    if gauge, ok := Gauges[name]; ok {
        log.Info(fmt.Sprintf("incrementing gauge '%s' %v", name, gauge))
        // generate labels for prometheus metric and check for errors
        labels, err := GenerateLabels(gaugeJson.Labels, "gauge", name)
        if err != nil {
            return err
        }
        gauge.With(labels).Inc()
        return nil
    }
    return ErrUnregisteredMetric
}

// function used to decrement a particular gauge value
func DecrementGauge(name string, gaugeJson GaugeJSON) error {
    if gauge, ok := Gauges[name]; ok {
        log.Info(fmt.Sprintf("decrementing gauge '%s' %v", name, gauge))
        // generate labels for prometheus metric and check for errors
        labels, err := GenerateLabels(gaugeJson.Labels, "gauge", name)
        if err != nil {
            return err
        }
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