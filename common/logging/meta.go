package logging

import (
	"go.uber.org/zap"
)

var version = "3.5.2"

func init() {
	Info("PasteMe Go Backend", zap.String("version", version))
}
