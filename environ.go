package main

import (
    "os"
    "fmt"
    "strconv"
    log "github.com/sirupsen/logrus"
)

var (
    LogLevels = map[string]log.Level{ "DEBUG": log.DebugLevel, "INFO": log.InfoLevel, "WARN": log.WarnLevel }
    ListenAddress string
    ListenPort int
    MaxBufferSize int
    ConfigFilePath string
)

// Function used to configure service settings
func ConfigureService() {
    // set log level by overriding environment variables
    LogLevelString := OverrideStringVariable("LOG_LEVEL", "DEBUG")
    if LogLevel, ok := LogLevels[LogLevelString]; ok {
        log.SetLevel(LogLevel)
    } else {
        log.Fatal(fmt.Sprintf("received invalid log level %s", LogLevelString))
    }
    // configure listen address and port from environment variables
    ListenAddress = OverrideStringVariable("LISTEN_ADDRESS", "0.0.0.0")
    ListenPort = OverrideIntegerVariable("LISTEN_PORT", 10091)

    MaxBufferSize = OverrideIntegerVariable("MAX_BUFFER_SIZE", 2048)
    ConfigFilePath = OverrideStringVariable("CONFIG_FILE_PATH", "./sample_config.json")
}

// Function used to override configuration variables with some
// value by defaulting from environment variables
func OverrideStringVariable(key string, DefaultValue string) string {
    value := os.Getenv(key)
    if len(value) > 0 {
        log.Info(fmt.Sprintf("overriding variable %v with value %v", key, value))
        return value
    } else {
        return DefaultValue
    }
}

// Function used to override configuration variables with some
// value by defaulting from environment variables
func OverrideIntegerVariable(key string, DefaultValue int) int {
    value := os.Getenv(key)
    if len(value) > 0 {
        result, err := strconv.Atoi(value)
        if err != nil {
            log.Fatal(fmt.Sprintf("cannot cast value '%v' to integer", result))
        }
        log.Info(fmt.Sprintf("overriding variable %v with value %v", key, result))
        return result
    } else {
        return DefaultValue
    }
}

// Function used to override configuration variables with some
// value by defaulting from environment variables
func OverrideFloatVariable(key string, DefaultValue float64) float64 {
    value := os.Getenv(key)
    if len(value) > 0 {
        result, err := strconv.ParseFloat(value, 64)
        if err != nil {
            log.Fatal(fmt.Sprintf("cannot cast value '%v' to float", result))
        }
        log.Info(fmt.Sprintf("overriding variable %v with value %v", key, result))
        return result
    } else {
        return DefaultValue
    }
}

// Function used to override configuration variables with some
// value by defaulting from environment variables
func OverrideBoolVariable(key string, DefaultValue bool) bool {
    value := os.Getenv(key)
    if len(value) > 0 {
        result, err := strconv.ParseBool(value)
        if err != nil {
            log.Fatal(fmt.Sprintf("cannot cast value '%v' to boolean", value))
        }
        log.Info(fmt.Sprintf("overriding variable %v with value %v", key, result))
        return result
    } else {
        return DefaultValue
    }
}

