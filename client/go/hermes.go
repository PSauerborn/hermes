package hermes_client

import (
    "fmt"
    "net"
    "errors"
    "encoding/json"
    log "github.com/sirupsen/logrus"
)

var (
    hermesHost *string
    hermesPort *int
    ErrHermesConnection = errors.New(fmt.Sprintf("cannot connect to hermes server"))
    ErrHermesPacketJSON = errors.New("unable to convert hermes udp packet to JSON format")
)

// define function used to send UDP packet to Hermes
// server. UDP Packets are converted to JSON before send
func sendUdpPacket(packet interface{}) error {
    // set default connection settings if settings not specified
    if hermesHost == nil || hermesPort == nil {
        setDefaultHermesConfig()
    }

    log.Debug(fmt.Sprintf("sending new udp packet %+v to hermes server", packet))
    // connect to hermes server and defer closing of connection
    conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", *hermesHost, *hermesPort))
    if err != nil {
        log.Error(fmt.Errorf("unable to connect to hermes server: %v", err))
        return errors.New(fmt.Sprintf("cannot connect to hermes server %s:%d", *hermesHost, *hermesPort))
    }
    defer conn.Close()

    // convert packet to JSON array
    bytes, err := json.Marshal(packet)
    if err != nil {
        log.Error(fmt.Errorf("unable to convert udp packet to JSON: %v", err))
        return ErrHermesPacketJSON
    }
    conn.Write(bytes)
    return nil
}

// function used to set default values for hermes configurations
func setDefaultHermesConfig() {
    host, port := "localhost", 7789
    hermesHost, hermesPort = &host, &port
}

// function used to set global hermes configuration
func SetHermesConfig(host string, port int) {
    log.Info(fmt.Sprintf("setting new Hermes config with %s:%d", host, port))
    hermesHost, hermesPort = &host, &port
}