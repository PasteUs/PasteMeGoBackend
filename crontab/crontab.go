package crontab

import (
	"github.com/robfig/cron/v3"
	"github.com/wonderivan/logger"
)
import "github.com/PasteUs/PasteMeGoBackend/model"

var entry cron.EntryID
var c *cron.Cron

// timePattern是类似Crontab模式的字符串,如"CRON_TZ=Asia/Shanghai 0 6 * * ?"代表6 a.m.运行
func StartClean(timePattern string) {
	c = cron.New()
	var err error
	entry, err = c.AddFunc(timePattern, func() {
		err := model.Clean()
		if err != nil {
			logger.Warn("Clean DataBase Failed,error:" + err.Error())
		}
	})
	if err != nil {
		logger.Warn("Create Crontab Failed,error:" + err.Error())

	}
	c.Start()
	logger.Info("Create Crontab Done.")

}

//func StopClean()  {
//c.Remove(entry)
//}
