package config

import (
    "encoding/json"
    "fmt"
    "github.com/PasteUs/PasteMeGoBackend/flag"
    "github.com/PasteUs/PasteMeGoBackend/meta"
    "github.com/PasteUs/PasteMeGoBackend/util/logging"
    "go.uber.org/zap"
    "io/ioutil"
    "os"
)

type database struct {
    Type     string `json:"type"`
    Username string `json:"username"`
    Password string `json:"password"`
    Server   string `json:"server"`
    Port     uint16 `json:"port"`
    Database string `json:"database"`
}

type Config struct {
    Version  string   `json:"version"`
    Address  string   `json:"address"`
    AdminUrl string   `json:"admin_url"` // PasteMe Admin's hostname
    Port     uint16   `json:"port"`
    Database database `json:"database"`
}

var config Config
var isInitialized bool

func init() {
    load(flag.Config)
    checkVersion(config.Version)
    setDefault()
}

func setDefault() {
}

func isInArray(item string, array []string) bool {
    for _, each := range array {
        if item == each {
            return true
        }
    }
    return false
}

func checkVersion(version string) {
    if version != meta.Version {

        if jsonBytes, err := json.Marshal(meta.ValidConfigVersion); err != nil {
            logging.Panic(err.Error())
        } else {
            if !isInArray(version, meta.ValidConfigVersion) {
                logging.Panic(fmt.Sprintf("valid config versions are %s, but \"%s\" was given", string(jsonBytes), version))
            }
        }
    }

}

func load(filename string) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        pwd, _ := os.Getwd()

        logging.Panic(
            "open file failed",
            zap.String("pwd", pwd),
            zap.String("err", err.Error()),
        )
    }

    err = json.Unmarshal(data, &config)
    if err != nil {
        logging.Panic(err.Error())
    }

    logging.Info(
        "config loaded",
        zap.String("config_file", filename),
        zap.ByteString("config", data),
    )

    isInitialized = true
}

func Get() Config {
    if !isInitialized {
        logging.Panic("Trying to use uninitialized config")
    }
    return config
}
