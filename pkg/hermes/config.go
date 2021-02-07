package hermes

import (
    "os"
    "fmt"
    "io/ioutil"
    "encoding/json"
    "errors"
    log "github.com/sirupsen/logrus"
)

var (
    ErrInvalidConfig = errors.New("Invalid hermes configuration")
)

// function used to generate HermesConfig instance from
// the local JSON configuration file. Additionally, if
// the listen address and port are not specified, the default
// values are assigned to the config
func LoadHermesConfig(path string) (HermesConfig, error) {
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
    return config, nil
}