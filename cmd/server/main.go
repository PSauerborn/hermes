package main

import (
    "strconv"

    log "github.com/sirupsen/logrus"

    "github.com/PSauerborn/hermes/pkg/hermes"
    "github.com/PSauerborn/hermes/pkg/utils"
)

var (
    // create map to house environment variables
    cfg = utils.NewConfigMapWithValues(
        map[string]string{
            "listen_port": "7832",
            "listen_address": "0.0.0.0",
        },
    )
)

func main() {
    log.SetLevel(log.DebugLevel)

    port, err := strconv.Atoi(cfg.Get("listen_port"))
    if err != nil {
        panic("received invalid listen port")
    }
    // start new instance of hermes server
    hermes.New("", cfg.Get("listen_address"), port).Listen()
}