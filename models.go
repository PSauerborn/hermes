package main


// struct used to define the global hermes configuration
// loaded for the local JSON file
type HermesConfig struct {
    ServiceName   string            `json:"service_name"`
    Gauges        []HermesGauge     `json:"gauges"`
    Counters      []HermesCounter   `json:"counters"`
    Histograms    []HermesHistogram `json:"histograms"`
    Summaries     []HermesSummary   `json:"summaries"`
    ListenAddress *string           `json:"listen_address"`
    ListenPort    *int              `json:"listen_port"`
}

// struct used to define a Gauge from the Hermes config
// used to create a prometheus gauge instance
type HermesGauge struct {
    Labels            []string `json:"labels"`
    MetricName        string   `json:"metric_name"`
    MetricDescription string   `json:"metric_description"`
}

// struct used to define a Counter from the Hermes config
// used to create a prometheus counter instance
type HermesCounter struct {
    Labels            []string `json:"labels"`
    MetricName        string   `json:"metric_name"`
    MetricDescription string   `json:"metric_description"`
}

// struct used to define a Counter from the Hermes config
// used to create a prometheus counter instance
type HermesHistogram struct {
    Labels            []string `json:"labels"`
    MetricName        string   `json:"metric_name"`
    MetricDescription string   `json:"metric_description"`
}

// struct used to define a Counter from the Hermes config
// used to create a prometheus counter instance
type HermesSummary struct {
    Labels            []string `json:"labels"`
    MetricName        string   `json:"metric_name"`
    MetricDescription string   `json:"metric_description"`
}

// struct used to define format of UDP packets
// sent from a hermes client
type HermesPayload struct {
    MetricName string      `json:"metric_name"`
    Payload    interface{} `json:"payload"`
}

// struct used to define JSON format of UDP packets
// for Gauges. Note that the operation field determines
// whether or not gauges are incremented, decremented
// or set with a particular value
type GaugeJSON struct {
    Labels    map[string]string `json:"labels"`
    Value     *float64          `json:"value"`
    Operation string            `json:"operation"`
}

// struct used to define JSON format of UDP packets for counters
type CounterJSON struct {
    Labels map[string]string `json:"labels"`
}

// struct used to define JSON format of UDP packets for counters
type HistogramJSON struct {
    Labels      map[string]string `json:"labels"`
    Observation float64           `json:"observation"`
}

type SummaryJSON struct {
    Labels      map[string]string `json:"labels"`
    Observation float64           `json:"observation"`
}