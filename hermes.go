package main

import (
    "fmt"
    "net"
    "time"
    "encoding/json"
    log "github.com/sirupsen/logrus"
)

var service *HermesServer

func main() {
    ConfigureService()
    // parse and load Hermes configuration from local JSON file
    config, err := LoadConfig(ConfigFilePath)
    if err != nil {
        log.Fatal(fmt.Errorf("error loading hermes configuration: %v", err))
    }
    log.Debug(fmt.Sprintf("starting new hermes server with configuration %+v", config))
    // create new hermes server instance and listen on go routine
    service = New(config)
    go ListenHermes(service)

    // begin listening to prometheus server on port 8080
    ListenPrometheus(config)
}

// function used to create new hermes service instance
func New(config HermesConfig) *HermesServer {
    addr := net.UDPAddr{IP: net.ParseIP(*config.ListenAddress), Port: *config.ListenPort}
    socket, err := net.ListenUDP("udp", &addr)
    if err != nil {
        log.Fatal(fmt.Errorf("unable to start new hermes server: %v", err))
    }
    return &HermesServer{Socket: socket, ListenAddress: &addr, Config: &config}
}

// function used to start listening on the specified UDP
// ports for JSON messages from a Hermes client. All incoming
// messages are read into a buffer and then converted to
// JSON format by the handler function
func ListenHermes(server *HermesServer) {
    // restart hermes socket if any panic issues arise during processing of messages
    defer func() {
        if r := recover(); r != nil {
            log.Warn(fmt.Sprintf("recovered paniced UDP interface: %+v", r))
            service.RestartServer(*server.Config)
        }
    }()
    // defer closing of connection
    defer service.Socket.Close()

    // create new buffer and serve messages
    buffer := make([]byte, MaxBufferSize)

    for {
        // read UDP packet payload into buffer
        n, remoteAddr, err := service.Socket.ReadFromUDP(buffer)
        log.Debug(fmt.Sprintf("processing new message from %+v", remoteAddr))
        if err != nil {
            log.Error(fmt.Errorf("unable to process UDP message: %v", err))
            continue
        }
        // handle UDP packet
        service.ProcessPayload(buffer[0:n])
    }
}

type HermesServer struct {
    Socket		  *net.UDPConn
    ListenAddress *net.UDPAddr
    Config 		  *HermesConfig
}

// function used to safely restart hermes server. the UDP connection
// is first closed via the socket connection. The connection is then
// re-established. If the re-creation of the socket fails, the go-routine
// will wait 10 seconds before attempting to re-open the connection
func(server HermesServer) RestartServer(config HermesConfig) {
    // close socket and restart
    service.Socket.Close()
    addr := net.UDPAddr{IP: net.ParseIP(*config.ListenAddress), Port: *config.ListenPort}
    for {
        socket, err := net.ListenUDP("udp", &addr)
        if err != nil {
            log.Error(fmt.Errorf("unable to start new hermes server: %v", err))
            time.Sleep(time.Second * 10)
        } else {
            service.Socket = socket
            break
        }
    }
    ListenHermes(&server)
}

// function used to process UDP packets sent over UDP interface.
// all packets are read into a buffer, and the contents of the
// buffer are then converted into JSON format. The metric name
// is sent with all JSON packets, which is then used to determine
// the type of metric that the JSON packet corresponds to (i.e.
// counter or gauge) and the payload is then processed depending on
// the type of metric
func(server HermesServer) ProcessPayload(packet []byte) {
    log.Debug(fmt.Sprintf("processing new hermes payload %s", string(packet)))
    var payload HermesPayload
    err := json.Unmarshal(packet, &payload)
    if err != nil {
        log.Error(fmt.Errorf("unable to parse udp packet to required JSON: %v", err))
        return
    }
    // determine metric type based on metric name
    metricType := GetMetricType(payload.MetricName)
    if metricType == nil {
        log.Error(fmt.Sprintf("cannot process metric %s. metric not registered", payload.MetricName))
        return
    }

    bytesPayload, _ := json.Marshal(payload.Payload)
    // process payload depending on metric type
    switch *metricType {
    case "counter":
        log.Debug(fmt.Sprintf("processing 'counter' payload %+v", payload.Payload))
        var counter CounterJSON
        err := json.Unmarshal(bytesPayload, &counter)
        if err != nil {
            log.Error(fmt.Sprintf("cannot process metric. invalid JSON"))
            return
        }
        IncrementCounter(payload.MetricName, counter)
    case "gauge":
        log.Debug(fmt.Sprintf("processing 'gauge' payload %+v", payload.Payload))
        var gauge GaugeJSON
        err := json.Unmarshal(bytesPayload, &gauge)
        if err != nil {
            log.Error(fmt.Sprintf("cannot process metric. invalid JSON"))
            return
        }
        ProcessGauge(payload.MetricName, gauge)
    }
}