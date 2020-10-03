package main

import (
    "fmt"
    "net"
    "encoding/json"
    log "github.com/sirupsen/logrus"
)

var service *HermesServer

// start new instance of hermes server
func main() {
    ConfigureService()
    // parse and load Hermes configuration from local JSON file
    config, err := LoadConfig(ConfigFilePath)
    if err != nil {
        log.Fatal(fmt.Errorf("error loading hermes configuration: %v", err))
    }
    log.Debug(fmt.Sprintf("starting new hermes server with configuration %+v", config))

    // begin listening to prometheus server on
    go ListenHermes(config)
    ListenPrometheus(config)
}

func ListenHermes(config HermesConfig) {
    // create new hermes instance
    service = New(config)
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
        // send message to heracles server to process message
        service.ProcessPayload(buffer[0:n])
    }
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

type HermesServer struct {
    Socket		  *net.UDPConn
    ListenAddress *net.UDPAddr
    Config 		  *HermesConfig
}

// define function used to process payload
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
        SetGauge(payload.MetricName, gauge)
    }
}