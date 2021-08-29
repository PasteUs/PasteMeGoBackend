package logging

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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

type loggerFunc func(string, ...interface{})

func getLogger(log func(string, ...zap.Field)) loggerFunc {
	return func(msg string, a ...interface{}) {
		log(msg, loggerPreprocess(a)...)
	}
}
