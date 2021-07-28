package config

import (
    "github.com/PasteUs/PasteMeGoBackend/util"
    "go.uber.org/zap"
)

var version = "3.4.0"
var validConfigVersion = []string{"3.3.0", ""}

func init() {
    util.Info("PasteMe Go Backend", zap.String("version", version))
}
