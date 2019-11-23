/*
@File: flag.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-07-25 08:36  Lucien     1.0         Init
*/
package flag

import (
	"flag"
	"fmt"
	"github.com/wonderivan/logger"
	"os"
	"strings"
)

var (
	version bool
	Config string
	Debug bool
	DataDir string
)

func init() {
	flag.BoolVar(&version, "version", false, "--version Print version information")
	flag.StringVar(&Config, "c", "./config.json", "-c <config file>")
	flag.BoolVar(&Debug, "debug", false, "--debug Using debug mode")
	flag.StringVar(&DataDir, "d", "./", "-d <data dir>")

	flag.Parse()

	validationCheck()
}

func validationCheck() {
	if !isDir(DataDir) {
		logger.Painc("%s is not a directory", DataDir)
	}

	if !strings.HasSuffix(DataDir, "/") {
		DataDir = DataDir + "/"
	}
}

func isDir(dataDir string) bool {
	if dir, err := os.Stat(dataDir); err == nil && dir != nil {
		return dir.IsDir()
	}
	return false
}

func Parse() bool { // return true for continue
	if version {
		fmt.Println("3.2.1")
		return false
	}
	return true
}
