package hermes_client

import (
    "fmt"

    log "github.com/sirupsen/logrus"
)


type HermesHistogramPacket struct {
    MetricName string                 `json:"metric_name"`
    Payload    HermesHistogramPayload `json:"payload"`
}

type HermesHistogramPayload struct {
    HistogramObservation float64           `json:"observation"`
    HistogramLabels      map[string]string `json:"labels"`
}

// function used to make an observation on a histogram metric
func(c *HermesClient) ObserveHistogram(metricName string, labels map[string]string, observation float64) {
    log.Debug(fmt.Sprintf("setting observation on histogram %s", metricName))
    packet := HermesHistogramPacket{
        MetricName: metricName,
        Payload: HermesHistogramPayload{
            HistogramLabels: labels,
            HistogramObservation: observation,
        },
    }
    c.SendUDPPacket(packet)
}