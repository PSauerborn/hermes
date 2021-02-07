package hermes_client

import (
    "fmt"

    log "github.com/sirupsen/logrus"
)


type HermesGaugePacket struct {
    MetricName string             `json:"metric_name"`
    Payload    HermesGaugePayload `json:"payload"`
}

type HermesGaugePayload struct {
    GaugeOperation string            `json:"operation"`
    GaugeValue	   *float64          `json:"value"`
    GaugeLabels    map[string]string `json:"labels"`
}

// function used to increment gauge value
func(c *HermesClient) IncrementGauge(metricName string, labels map[string]string) {
    log.Debug(fmt.Sprintf("incrementing gauge %s", metricName))
    // generate new packet and send via UDP socket
    packet := HermesGaugePacket{
        MetricName: metricName,
        Payload: HermesGaugePayload{GaugeOperation: "increment", GaugeLabels: labels},
    }
    c.SendUDPPacket(packet)
}

// function used to decrement gauge value
func(c *HermesClient) DecrementGauge(metricName string, labels map[string]string) {
    log.Debug(fmt.Sprintf("decrementing gauge %s", metricName))
    // generate new packet and send via UDP socket
    packet := HermesGaugePacket{
        MetricName: metricName,
        Payload: HermesGaugePayload{GaugeOperation: "decrement", GaugeLabels: labels},
    }
    c.SendUDPPacket(packet)
}

// function used to set value ot gauge to a user defined float value
func(c *HermesClient) SetGauge(metricName string, labels map[string]string, gaugeValue float64) {
    log.Debug(fmt.Sprintf("setting gauge %s with value %f", metricName, gaugeValue))
    // generate new packet and send via UDP socket
    packet := HermesGaugePacket{
        MetricName: metricName,
        Payload: HermesGaugePayload{
            GaugeOperation: "set",
            GaugeValue: &gaugeValue,
            GaugeLabels: labels,
        },
    }
    c.SendUDPPacket(packet)
}