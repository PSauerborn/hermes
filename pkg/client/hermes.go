package hermes_client

import (
    "fmt"
    "net"
    "errors"
    "encoding/json"

    log "github.com/sirupsen/logrus"
)

var (
    // define custom errors
    ErrHermesConnection = errors.New("Cannot connect to hermes server")
    ErrHermesPacketJSON = errors.New("Unable to convert hermes udp packet to JSON format")
)

// struct used to container hermes client details
type HermesClient struct {
    HermesHost string
    HermesPort int
}

// function used to generate new hermes client
func New(host string, port int) *HermesClient {
    return &HermesClient{
        HermesHost: host,
        HermesPort: port,
    }
}

// define function used to send UDP packet to Hermes
// server. UDP Packets are converted to JSON before send
func(c *HermesClient) SendUDPPacket(packet interface{}) error {
    log.Debug(fmt.Sprintf("sending new udp packet %+v to hermes server", packet))
    // connect to hermes server and defer closing of connection
    conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", c.HermesHost, c.HermesPort))
    if err != nil {
        log.Error(fmt.Errorf("unable to connect to hermes server: %v", err))
        return errors.New(fmt.Sprintf("cannot connect to hermes server %s:%d", c.HermesHost, c.HermesPort))
    }
    defer conn.Close()

    // convert JSON packet into bytes array
    bytes, err := json.Marshal(packet)
    if err != nil {
        log.Error(fmt.Errorf("unable to convert udp packet to JSON: %v", err))
        return ErrHermesPacketJSON
    }
    // write hermes packet over UDP socket
    conn.Write(bytes)
    return nil
}