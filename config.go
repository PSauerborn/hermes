package main

import (
    "os"
    "fmt"
    "io/ioutil"
    "encoding/json"
    "errors"
    log "github.com/sirupsen/logrus"
)

var (
    ErrInvalidConfig = errors.New("invalid hermes configuration")
)


// function used to validate the Hermes configuration instance
func ValidateConfig(config HermesConfig) error {
    return nil
}

// function used to load config from local file path
func LoadConfig(path string) (HermesConfig, error) {
    var config HermesConfig

    configFile, err := os.Open(path)
    if err != nil {
        log.Error(fmt.Errorf("cannot load local JSON configuration: %v", err))
        return config, ErrInvalidConfig
    }
    defer configFile.Close()
    // convert to bytes using ioutil
    bytesJson, err := ioutil.ReadAll(configFile)
    if err != nil {
        log.Error(fmt.Errorf("cannot load local JSON configuration: %v", err))
        return config, ErrInvalidConfig
    }
    // cast to JSON format and return
    err = json.Unmarshal(bytesJson, &config)
    if err != nil {
        log.Error(fmt.Errorf("cannot load local JSON configuration: %v", err))
        return config, ErrInvalidConfig
    }
    // set default listen address if not specified
    if config.ListenAddress == nil {
        address := "0.0.0.0"; config.ListenAddress = &address
    }
    // set default listen port if not specified
    if config.ListenPort == nil {
        port := 7789; config.ListenPort = &port
    }
    return config, nil
}