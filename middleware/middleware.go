package middleware

import (
	"github.com/PasteUs/PasteMeGoBackend/config"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/wonderivan/logger"
	"os"
	"path"
	"time"
)

func init()  {

	LogFilePath := config.Data.LogFilePath
	LogFileName := config.Data.LogFileName
	/*日志文件*/
	fileName := path.Join(LogFilePath,LogFileName)

	f,err := os.Create( fileName )
	defer f.Close()

	if err !=nil {
		logger.Error("Can't create file,Please check for errors in config.json")
	}
}
/*记录到文件*/
func LoggerToFile() gin.HandlerFunc{

	LogFilePath := config.Data.LogFilePath
	LogFileName := config.Data.LogFileName
	/*日志文件*/
	fileName := path.Join(LogFilePath,LogFileName)

	src,err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if(err!=nil){
		logger.Error("can't find file",err)
	}
	logrusInstance := logrus.New()
	logrusInstance.Out = src

	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName + ".%Y%m%d.log",

		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),

		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),

		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat:"2006-01-02 15:04:05",
	})

	// 新增 Hook
	logrusInstance.AddHook(lfHook)

	return func(context *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		context.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := context.Request.Method

		// 请求路由
		reqUri := context.Request.RequestURI

		// 状态码
		statusCode := context.Writer.Status()

		// 请求IP
		clientIP := context.ClientIP()

		// 日志格式
		logrusInstance.WithFields(logrus.Fields{
			"status_code"  : statusCode,
			"latency_time" : latencyTime,
			"client_ip"    : clientIP,
			"req_method"   : reqMethod,
			"req_uri"      : reqUri,
		}).Info()

	}
}