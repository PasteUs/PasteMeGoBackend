package config

import (
	"encoding/json"
	"github.com/PasteUs/PasteMeGoBackend/flag"
	"github.com/PasteUs/PasteMeGoBackend/logging"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

type Database struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Port     uint16 `json:"port"`
	Database string `json:"database"`
}

type config struct {
	Address  string   `json:"address"`
	Port     uint16   `json:"port"`
	Secret   string   `json:"secret"`
	LogFile  string   `json:"log_file"`
	Database Database `json:"database"`
}

var Config config

func init() {
	load(flag.Config)
}

func exportConfig(filename string, c config) {
	if flag.Debug {
		logging.Info(
			"Config loaded",
			zap.String("config_file", filename),
			zap.String("address", c.Address),
			zap.String("log_file", c.LogFile),
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

	err = json.Unmarshal(data, &Config)
	if err != nil {
		logging.Panic("parse Config failed", zap.String("err", err.Error()))
	}

	exportConfig(filename, Config)
}
