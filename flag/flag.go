package flag

import (
    "flag"
    "github.com/wonderivan/logger"
    "os"
    "strings"
    "testing"
)

var (
    Config  string
    Debug   bool
    DataDir string
)

func init() {
    flag.StringVar(&Config, "c", "config.json", "-c <config file>")
    flag.BoolVar(&Debug, "debug", false, "--debug Using debug mode")
    flag.StringVar(&DataDir, "d", "./", "-d <data dir>")

    testing.Init()
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
