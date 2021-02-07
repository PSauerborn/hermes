package hermes

import (
    "fmt"
    "net"
    "time"
    "encoding/json"
    log "github.com/sirupsen/logrus"
)

type HermesServer struct {
    // UDP socket to listen for packets
    Socket		  *net.UDPConn
    ListenAddress *net.UDPAddr

    // hermes config containing data about metrics
    Config 		  HermesConfig
}

// function used to create new hermes service instance
func New(configPath, listenAddress string, listenPort int) *HermesServer {
    // load hermes configration from local file
    cfg, err := LoadHermesConfig(configPath)
    if err != nil {
        panic(fmt.Errorf("unable to load hermes config from path: %s: %+v", configPath,
            err))
    }
    // generate new UDP address instance and socket to listen on
    addr := net.UDPAddr{IP: net.ParseIP(listenAddress), Port: listenPort}
    socket, err := net.ListenUDP("udp", &addr)
    if err != nil {
        log.Fatal(fmt.Errorf("unable to start new hermes server: %v", err))
    }
    return &HermesServer{Socket: socket, ListenAddress: &addr, Config: cfg}
}

// function used to start listening on the specified UDP
// ports for JSON messages from a Hermes client. All incoming
// messages are read into a buffer and then converted to
// JSON format by the handler function
func(server *HermesServer) Listen() {
    log.Info(fmt.Sprintf("starting new UDP interface at %+v...", server.ListenAddress))
    // restart hermes socket if any panic issues arise during processing of messages
    defer func() {
        if r := recover(); r != nil {
            log.Warn(fmt.Sprintf("recovered paniced UDP interface: %+v", r))
            server.RestartServerGracefully()
        }
    }()
    // defer closing of connection
    defer server.Socket.Close()
    // start HTTP Prometheus server on goroutine
    go ListenPrometheus(server.Config, 8080)

    // create new buffer and serve messages
    buffer := make([]byte, 2048)
    for {
        // read UDP packet payload into buffer
        n, remoteAddr, err := server.Socket.ReadFromUDP(buffer)
        log.Debug(fmt.Sprintf("processing new message from %+v", remoteAddr))
        if err != nil {
            log.Error(fmt.Errorf("unable to process UDP message: %v", err))
            continue
        }
        // handle UDP packet
        server.ProcessPayload(buffer[0:n])
    }
}

// function used to safely restart hermes server. the UDP connection
// is first closed via the socket connection. The connection is then
// re-established. If the re-creation of the socket fails, the go-routine
// will wait 10 seconds before attempting to re-open the connection
func(server *HermesServer) RestartServerGracefully() {
    // close socket and restart
    server.Socket.Close()
    for {
        // attempt to recreate socket connection (active connections can
        // take a while to stop properly) and restart
        socket, err := net.ListenUDP("udp", server.ListenAddress)
        if err != nil {
            log.Error(fmt.Errorf("unable to start new hermes server: %v", err))
            time.Sleep(time.Second * 10)
        } else {
            server.Socket = socket
            break
        }
    }
    // start listening on socket once connection has been setup
    server.Listen()
}

// function used to process UDP packets sent over UDP interface.
// all packets are read into a buffer, and the contents of the
// buffer are then converted into JSON format. The metric name
// is sent with all JSON packets, which is then used to determine
// the type of metric that the JSON packet corresponds to (i.e.
// counter or gauge) and the payload is then processed depending on
// the type of metric
func(server *HermesServer) ProcessPayload(packet []byte) {
    log.Debug(fmt.Sprintf("processing new hermes payload %s", string(packet)))
    var payload HermesPayload
    err := json.Unmarshal(packet, &payload)
    if err != nil {
        log.Error(fmt.Errorf("unable to parse udp packet to required JSON format: %v", err))
        return
    }
    // determine metric type based on metric name from local mappings of metrics
    metricType, err := GetMetricType(payload.MetricName)
    if err != nil {
        log.Error(fmt.Sprintf("cannot process metric %s: metric not registered", payload.MetricName))
        return
    }

    bytesPayload, _ := json.Marshal(payload.Payload)
    // process payload depending on metric type
    switch metricType {

    // process counter metrics
    case "counter":
        log.Debug(fmt.Sprintf("processing 'counter' payload %+v", payload.Payload))
        var counter CounterJSON
        err := json.Unmarshal(bytesPayload, &counter)
        if err != nil {
            log.Error(fmt.Sprintf("cannot process 'counter' metric. invalid JSON"))
            return
        }
        IncrementCounter(payload.MetricName, counter)

    // process gauge metrics
    case "gauge":
        log.Debug(fmt.Sprintf("processing 'gauge' payload %+v", payload.Payload))
        var gauge GaugeJSON
        err := json.Unmarshal(bytesPayload, &gauge)
        if err != nil {
            log.Error(fmt.Sprintf("cannot process 'gauge' metric. invalid JSON"))
            return
        }
        ProcessGauge(payload.MetricName, gauge)

    // process histogram metrics
    case "histogram":
        log.Debug(fmt.Sprintf("processing 'histogram' payload %+v", payload.Payload))
        var histogram HistogramJSON
        err := json.Unmarshal(bytesPayload, &histogram)
        if err != nil {
            log.Error(fmt.Sprintf("cannot process 'histogram' metric. invalid JSON"))
            return
        }
        ObserveHistogram(payload.MetricName, histogram)

    // process summary metrics
    case "summary":
        log.Debug(fmt.Sprintf("processing 'summary' payload %+v", payload.Payload))
        var summary SummaryJSON
        err := json.Unmarshal(bytesPayload, &summary)
        if err != nil {
            log.Error(fmt.Sprintf("cannot process 'summary' metric. invalid JSON"))
            return
        }
        ObserveSummary(payload.MetricName, summary)
    }
}