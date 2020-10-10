package hermes_client

import (
    "fmt"
    log "github.com/sirupsen/logrus"
)


type HermesSummaryPacket struct {
    MetricName string               `json:"metric_name"`
    Payload    HermesSummaryPayload `json:"payload"`
}

type HermesSummaryPayload struct {
    SummaryObservation float64           `json:"observation"`
    SummaryLabels      map[string]string `json:"labels"`
}

// function used to make an observation on a summary metric
func ObserveSummary(metricName string, labels map[string]string, observation float64) {
    log.Debug(fmt.Sprintf("setting observation on histogram %s", metricName))
    packet := HermesSummaryPacket{
        MetricName: metricName,
        Payload: HermesSummaryPayload{
            SummaryLabels: labels,
            SummaryObservation: observation,
        },
    }
    sendUdpPacket(packet)
}