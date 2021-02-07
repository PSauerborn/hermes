package hermes_client

import (
    "fmt"

    log "github.com/sirupsen/logrus"
)


type HermesCounterPacket struct {
    MetricName string               `json:"metric_name"`
    Payload    HermesCounterPayload `json:"payload"`
}

type HermesCounterPayload struct {
    CounterLabels map[string]string `json:"labels"`
}

// function used to increment counter value
func IncrementCounter(metricName string, labels map[string]string) {
    log.Debug(fmt.Sprintf("incrementing counter %s", metricName))
    packet := HermesCounterPacket{
        MetricName: metricName,
        Payload: HermesCounterPayload{
            CounterLabels: labels,
        },
    }
    sendUdpPacket(packet)
}