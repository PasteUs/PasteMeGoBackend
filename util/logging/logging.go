package logging

import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

var logger *zap.Logger

func init() {
    config := zap.NewProductionConfig()

    config.OutputPaths = append(config.OutputPaths, "pasteme.log")

    lgr, _ := config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
    logger = lgr
}

func exportField(requests *gin.Context) []zap.Field {
    var result []zap.Field
    result = append(result, zap.String("ip", requests.ClientIP()))
    return result
}

func loggerPreprocess(fields []interface{}) []zap.Field {
    var result []zap.Field

    if len(fields) != 0 {
        var beginIndex = 0
        switch fields[0].(type) {
        case *gin.Context:
            context := fields[0].(*gin.Context)
            zapField := exportField(context)
            result = append(result, zapField...)
            beginIndex = 1
        }
        for _, each := range fields[beginIndex:] {
            result = append(result, each.(zap.Field))
        }
    }

    return result
}

func Info(msg string, a ...interface{}) {
    fields := loggerPreprocess(a)
    logger.Info(msg, fields...)
}

func Warn(msg string, a ...interface{}) {
    fields := loggerPreprocess(a)
    logger.Warn(msg, fields...)
}

func Panic(msg string, a ...interface{}) {
    fields := loggerPreprocess(a)
    logger.Panic(msg, fields...)
}
