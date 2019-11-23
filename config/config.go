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
	"github.com/wonderivan/logger"
	"io/ioutil"
)

type database struct {
	Type string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
	Server string `json:"server"`
	Port uint16 `json:"port"`
	Database string `json:"database"`
}

type Config struct {
	Address string `json:"address"`
	Port uint16 `json:"port"`
	Database database `json:"database"`
}

var config Config
var isInitialized bool

func Load(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Fatal(err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		logger.Fatal(err)
	}

	isInitialized = true
}

func Get() Config {
	if !isInitialized {
		logger.Fatal("Trying to use uninitialized config")
	}
	return config
}
