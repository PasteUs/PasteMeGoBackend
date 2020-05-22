/*
@File: config.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-07-25 01:33  Lucien     1.0         Init
*/
package config

import (
	"encoding/json"
	"github.com/PasteUs/PasteMeGoBackend/flag"
	"github.com/wonderivan/logger"
	"io/ioutil"
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
	Address  string   `json:"address"`
	Port     uint16   `json:"port"`
	Database database `json:"database"`
}

var config Config
var isInitialized bool

func init() {
	load(flag.Config)
}

func load(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Painc(err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		logger.Painc(err)
	}

	logger.Info("Load from %s\nconfig = %s", filename, data)

	isInitialized = true
}

func Get() Config {
	if !isInitialized {
		logger.Painc("Trying to use uninitialized config")
	}
	return config
}
