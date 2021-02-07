package main

import (
    "github.com/PSauerborn/hermes/pkg/client"
)

func main() {
    // set configuration for Hermes client and increment sample counter
    client := hermes_client.New("192.168.99.100", 7789)

    // increment counter metric over hermes server
    client.IncrementCounter("sample_counter",
        map[string]string{"label_1": "test-label"})

    // increment gauge
    client.IncrementGauge("sample_gauge",
        map[string]string{"label_1": "test-label", "label_2": "test-label-2"})

    // decrement gauge
    client.DecrementGauge("sample_gauge",
        map[string]string{"label_1": "test-label", "label_2": "test-label-2"})

    // set value of gauge
    client.SetGauge("sample_gauge",
        map[string]string{"label_1": "test-label", "label_2": "test-label-random"}, 54)

    // make observation on histogram
    client.ObserveHistogram("sample_histogram",
        map[string]string{"label_1": "test-label"}, 1233)

    // make observation on summary
    client.ObserveSummary("sample_summary",
        map[string]string{"label_1": "test-label"}, 5)
}