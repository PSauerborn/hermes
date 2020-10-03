package main


type HermesConfig struct {
    ServiceName   string          `json:"service_name"`
    Gauges        []HermesGauge   `json:"gauges"`
    Counters      []HermesCounter `json:"counters"`
    ListenAddress *string         `json:"listen_address"`
    ListenPort    *int            `json:"listen_port"`
}

type HermesPayload struct {
    MetricName string      `json:"metric_name"`
    Payload    interface{} `json:"payload"`
}

type GaugeJSON struct {
    Labels map[string]string `json:"labels"`
    Value  float64           `json:"value"`
}

type CounterJSON struct {
    Labels map[string]string `json:"labels"`
}

type HermesGauge struct {
    Labels            []string `json:"labels"`
    MetricName        string   `json:"metric_name"`
    MetricDescription string   `json:"metric_description"`
}

type HermesCounter struct {
    Labels            []string `json:"labels"`
    MetricName        string   `json:"metric_name"`
    MetricDescription string   `json:"metric_description"`
}