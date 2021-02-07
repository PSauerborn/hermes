package hermes

import (
    "fmt"

    "github.com/prometheus/client_golang/prometheus"
    log "github.com/sirupsen/logrus"
)

// function used to increment a particular counter
func IncrementCounter(name string, counterJson CounterJSON) error {
    if counter, ok := Counters[name]; ok {
        log.Info(fmt.Sprintf("incrementing counter '%s' %v", name, counter))
        // generate labels for prometheus metric and check for errors
        labels, err := GenerateLabels(counterJson.Labels, "counter", name)
        if err != nil {
            return err
        }
        counter.With(labels).Inc()
        return nil
    }
    return ErrUnregisteredMetric
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