package config

import (
	"encoding/json"
	"github.com/PasteUs/PasteMeGoBackend/flag"
	"github.com/PasteUs/PasteMeGoBackend/logging"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"sync"
)

type Config struct {
	Version  string `json:"version"`
	Address  string `json:"address"`
	AdminUrl string `json:"admin_url"` // PasteMe Admin's hostname
	Port     uint16 `json:"port"`
	Database struct {
		Type     string `json:"type"`
		Username string `json:"username"`
		Password string `json:"password"`
		Server   string `json:"server"`
		Port     uint16 `json:"port"`
		Database string `json:"database"`
	} `json:"database"`
}

var (
	config Config
	once   sync.Once
)

func Init() {
	load(flag.GetArgv().Config)
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

func checkVersion(v string) {
	if v != version {
		if !isInArray(v, validConfigVersion) {
			logging.Panic(
				"invalid config version",
				zap.Strings("valid_config_version_list", validConfigVersion),
				zap.String("config_version", v),
			)
		}
	}
}

func exportConfig(filename string, c Config) {
	if flag.GetArgv().Debug {
		logging.Info(
			"config loaded",
			zap.String("config_file", filename),
			zap.String("config_version", c.Version),
			zap.String("address", c.Address),
			zap.String("admin_url", c.AdminUrl),
			zap.Uint16("port", c.Port),
		)

		if c.Database.Type == "mysql" {
			logging.Info(
				"database",
				zap.String("type", c.Database.Type),
				zap.String("username", c.Database.Username),
				zap.String("password", c.Database.Password),
				zap.String("server", c.Database.Server),
				zap.Uint16("port", c.Database.Port),
				zap.String("database", c.Database.Database),
			)
		} else {
			logging.Info(
				"database",
				zap.String("type", c.Database.Type),
			)
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
		logging.Panic("parse config failed", zap.String("err", err.Error()))
	}

	exportConfig(filename, config)
}

func Get() Config {
	once.Do(Init)
	return config
}
