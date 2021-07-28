package flag

import (
    "flag"
    "github.com/PasteUs/PasteMeGoBackend/util"
    "go.uber.org/zap"
    "os"
    "strings"
    "sync"
)

type Argv struct {
    Config  string
    Debug   bool
    DataDir string
}

var (
    argv Argv
    once sync.Once
)

func init() {
    flag.StringVar(&argv.Config, "c", "config.json", "-c <config file>")
    flag.BoolVar(&argv.Debug, "debug", false, "--debug Using debug mode")
    flag.StringVar(&argv.DataDir, "d", "./", "-d <data dir>")
}

func Init() {
    flag.Parse()
    validationCheck(argv.DataDir)
}

func GetArgv() Argv {
    once.Do(Init)
    return argv
}

func validationCheck(dataDir string) {
    if !isDir(dataDir) {
        util.Panic("not a directory", zap.String("data_dir", dataDir))
    }

    if !strings.HasSuffix(dataDir, "/") {
        dataDir = dataDir + "/"
    }
}

func isDir(dataDir string) bool {
    if dir, err := os.Stat(dataDir); err == nil && dir != nil {
        return dir.IsDir()
    }
    return false
}
