package logging

import (
	"github.com/PasteUs/PasteMeGoBackend/common/config"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	Debug  loggerFunc
	Info   loggerFunc
	Warn   loggerFunc
	Error  loggerFunc
	Fatal  loggerFunc
	Panic  loggerFunc
)

func init() {
	zapConfig := zap.NewProductionConfig()
	if config.Config.LogFile != "" {
		zapConfig.OutputPaths = append(zapConfig.OutputPaths, config.Config.LogFile)
	}

	lgr, _ := zapConfig.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	logger = lgr

	Debug = getLogger(logger.Debug)
	Info = getLogger(logger.Info)
	Warn = getLogger(logger.Warn)
	Error = getLogger(logger.Error)
	Fatal = getLogger(logger.Fatal)
	Panic = getLogger(logger.Panic)
}
