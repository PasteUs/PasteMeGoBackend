package config

import (
	"github.com/PasteUs/PasteMeGoBackend/common/logging"
	"go.uber.org/zap"
)

var version = "3.5.2"

func init() {
	logging.Info("PasteMe Go Backend", zap.String("version", version))
}
