package config

import (
	"encoding/json"
	"fmt"
	"github.com/PasteUs/PasteMeGoBackend/common/flag"
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

func load(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		pwd, _ := os.Getwd()
		panic(fmt.Sprintf("open %s under %s failed, error = \"%s\"\n", filename, pwd, err.Error()))
	}

	if flag.Debug {
		fmt.Printf("read %s done, got %s\n", filename, string(data))
	}

	err = json.Unmarshal(data, &Config)
	if err != nil {
		panic(fmt.Sprintf("parse Config failed, error = \"%s\"\n", err.Error()))
	}
}
