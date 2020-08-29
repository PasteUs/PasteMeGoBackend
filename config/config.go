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
	"github.com/PasteUs/PasteMeGoBackend/meta"
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
			logger.Painc(err)
		} else {
			if !isInArray(version, meta.ValidConfigVersion) {
				logger.Painc("Valid config versions are %s, but \"%s\" was given", string(jsonBytes), version)
			}
		}
	}

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
