package main

import (
    "strings"
    "strconv"

    log "github.com/sirupsen/logrus"

    "github.com/PSauerborn/hermes/pkg/hermes"
    "github.com/PSauerborn/hermes/pkg/utils"
)

var (
    // create map to house environment variables
    cfg = utils.NewConfigMapWithValues(
        map[string]string{
            "listen_port": "7789",
            "listen_address": "0.0.0.0",
            "hermes_config_path" : "/etc/hermes/config.json",
            "log_level": "INFO",
        },
    )
)

// function to set log level from environment variables
func SetLogLevel() {
    level := cfg.Get("log_level")
    switch strings.ToUpper(level) {
    case "DEBUG":
        log.SetLevel(log.DebugLevel)
    case "INFO":
        log.SetLevel(log.InfoLevel)
    case "WARN":
        log.SetLevel(log.WarnLevel)
    case "ERROR":
        log.SetLevel(log.ErrorLevel)
    }
}

func main() {
    // set log level for server
    SetLogLevel()

    port, err := strconv.Atoi(cfg.Get("listen_port"))
    if err != nil {
        panic("received invalid listen port")
    }
    // start new instance of hermes server
    hermes.New(cfg.Get("hermes_config_path"), cfg.Get("listen_address"),
        port).Listen()
}