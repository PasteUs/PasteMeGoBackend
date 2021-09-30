package logging

import (
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
	config := zap.NewProductionConfig()

	config.OutputPaths = append(config.OutputPaths, "pasteme.log")

	lgr, _ := config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	logger = lgr

	Debug = getLogger(logger.Debug)
	Info = getLogger(logger.Info)
	Warn = getLogger(logger.Warn)
	Error = getLogger(logger.Error)
	Fatal = getLogger(logger.Fatal)
	Panic = getLogger(logger.Panic)
}
