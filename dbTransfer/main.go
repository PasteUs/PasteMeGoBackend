/*
@File: main.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

Transform database from version 2.x to version 3.x

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-07-01 23:39  Lucien     1.0         Init
*/
package main

import (
	"fmt"
	"github.com/PasteUs/PasteMeGoBackend/dbTransfer/database"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wonderivan/logger"
	"time"
)

func main() {
	database.FixAutoIncrement()

	start := time.Now()
	permCount := database.TransPermanent()
	timePerm := time.Since(start)
	logger.Info("Permanent finished: ", timePerm)

	start = time.Now()
	tempCount := database.TransTemporary()
	timeTemp := time.Since(start)
	logger.Info("Temporary finished: ", timeTemp)

	logger.Info("=====================================")
	logger.Info(fmt.Sprintf("%d records total, cost: ", permCount+tempCount), timePerm+timeTemp)
	logger.Info(fmt.Sprintf("%d permanents cost: ", permCount), timePerm)
	logger.Info(fmt.Sprintf("%d temporaries cost: ", tempCount), timeTemp)
}
