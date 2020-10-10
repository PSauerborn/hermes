package main

import (
    "github.com/PSauerborn/hermes/client/go"
)

func main() {
    // set configuration for Hermes client and increment sample counter
    hermes_client.SetHermesConfig("localhost", 7789)

    // increment counter metric over hermes server
    hermes_client.IncrementCounter("sample_counter", map[string]string{"label_1": "test-label"})

    // set gauge values over hermes server
    hermes_client.IncrementGauge("sample_gauge", map[string]string{"label_1": "test-label", "label_2": "test-label-2"})
    hermes_client.DecrementGauge("sample_gauge", map[string]string{"label_1": "test-label", "label_2": "test-label-2"})
    hermes_client.SetGauge("sample_gauge", map[string]string{"label_1": "test-label", "label_2": "test-label-random"}, 54)

    // make observation on histogram
    hermes_client.ObserveHistogram("sample_histogram", map[string]string{"label_1": "test-label"}, 1233)

    // make observation on summary
    hermes_client.ObserveSummary("sample_summary", map[string]string{"label_1": "test-label"}, 5)
}