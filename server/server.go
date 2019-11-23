/*
@File: server.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-06-21 08:37  Lucien     1.0         Init
*/
package server

import (
	"fmt"
	"github.com/PasteUs/PasteMeGoBackend/model"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"os"
)

var router *gin.Engine

func init() {
	if os.Getenv("PASTEMED_RUNTIME") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	router = gin.Default()
	router.GET("/", beat) // 心跳检测
	// 访问未加密的 Paste，token 为 <Paste ID>
	// 访问加密的 Paste，token 为 <Paste ID>,<Password>
	router.GET("/:token", query)
	router.POST("/", permanentCreator) // 创建一个永久的 Paste, key 是自增键
	router.POST("/once", readOnceCreator) // 创建一个阅后即焚的 Paste, key 是随机的
	router.PUT("/:key", temporaryCreator) // 创建一个阅后即焚的 Paste, key 是指定的
	router.NoRoute(notFoundHandler)
}

func Run(address string, port uint16) {
	model.Init()
	if err := router.Run(fmt.Sprintf("%s:%d", address, port)); err != nil {
		logger.Fatal("Run server failed: " + err.Error())
	}
}
