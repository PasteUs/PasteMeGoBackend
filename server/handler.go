/*
@File: handler.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-06-23 16:02  Lucien     1.0         None
*/
package server

import (
	"fmt"
	"github.com/LucienShui/PasteMeBackend/model"
	"github.com/LucienShui/PasteMeBackend/util"
	"github.com/LucienShui/PasteMeBackend/util/convert"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
	"strings"
)

func permanent(requests *gin.Context) {
	IP := requests.ClientIP()
	paste := model.Permanent{
		ClientIP: IP,
	}
	if err := requests.ShouldBindJSON(&paste); err != nil {
		logger.Error(util.LoggerInfo(IP, "Bind failed: "+err.Error()))
		requests.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "bind failed",
			"error":   err.Error(),
		})
	} else {
		if err := paste.Save(); err != nil {
			logger.Error(util.LoggerInfo(IP, "Save failed: "+err.Error()))
			requests.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "save failed",
				"error":   err.Error(),
			})
		} else {
			logger.Info(util.LoggerInfo(IP, "Create an permanent paste with key: "+convert.Uint2string(paste.Key)))
			requests.JSON(http.StatusCreated, gin.H{
				"status": http.StatusCreated,
				"key":    paste.Key,
			})
		}
	}
}

func temporary(requests *gin.Context) {
	IP, key := requests.ClientIP(), requests.Param("key")
	key = strings.ToLower(key)
	table, err := util.ValidChecker(key)
	if err != nil {
		if err.Error() == "wrong length" {
			logger.Warn(util.LoggerInfo(IP, "Trying to create temporary paste with key: "+key))
			requests.JSON(http.StatusOK, gin.H{
				"status":  http.StatusBadRequest,
				"error":   err.Error(),
				"message": "key's length should at least 3 and at most 8",
			})
		} else {
			logger.Warn(util.LoggerInfo(IP, "Trying to create temporary paste with key: "+key))
			requests.JSON(http.StatusOK, gin.H{
				"status":  http.StatusBadRequest,
				"error":   err.Error(),
				"message": "temporary key should only contains digits and lowercase letters, at least one alpha is required",
			})
		}
	} else {
		if table != "temporary" {
			logger.Warn(util.LoggerInfo(IP, "Trying to create temporary paste with key: "+key))
			requests.JSON(http.StatusOK, gin.H{
				"status":  http.StatusBadRequest,
				"error":   "wrong key type",
				"message": "temporary key should only contains digits and lowercase letters, at least one alpha is required",
			})
		} else {
			paste := model.Temporary{
				Key:      key,
				ClientIP: requests.ClientIP(),
			}
			if err := requests.ShouldBindJSON(&paste); err != nil {
				logger.Error(util.LoggerInfo(IP, "Bind failed: "+err.Error()))
				requests.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"error":   err.Error(),
					"message": "bind failed",
				})
			} else {
				if err := paste.Save(); err != nil {
					logger.Error(util.LoggerInfo(IP, "Save failed: "+err.Error()))
					requests.JSON(http.StatusInternalServerError, gin.H{
						"status":  http.StatusInternalServerError,
						"error":   err.Error(),
						"message": "save failed",
					})
				} else {
					logger.Info(util.LoggerInfo(IP, "Create an temporary paste with key: "+paste.Key))
					requests.JSON(http.StatusCreated, gin.H{
						"status": http.StatusCreated,
						"key":    paste.Key,
					})
				}
			}
		}
	}
}

func readOnce(requests *gin.Context) {
	IP := requests.ClientIP()
	paste := model.Temporary{
		Key:      util.Generator(),
		ClientIP: IP,
	}
	if err := requests.ShouldBindJSON(&paste); err != nil {
		logger.Error(util.LoggerInfo(IP, "Bind failed: "+err.Error()))
		requests.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"error":   err.Error(),
			"message": "bind failed",
		})
	} else {
		if err := paste.Save(); err != nil {
			logger.Error(util.LoggerInfo(IP, "Save failed: "+err.Error()))
			requests.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"error":   err.Error(),
				"message": "save failed",
			})
		} else {
			logger.Info(util.LoggerInfo(IP, "Create an once paste with key: "+paste.Key))
			requests.JSON(http.StatusCreated, gin.H{
				"status": http.StatusCreated,
				"key":    paste.Key,
			})
		}
	}
}

func get(requests *gin.Context) {
	IP, token := requests.ClientIP(), requests.Param("token")
	if token == "" { // empty token
		requests.JSON(http.StatusOK, gin.H{
			"status":  http.StatusBadRequest,
			"error":   "empty token",
			"message": "wrong params",
		})
	} else {
		key, password := util.Parse(token)
		key = strings.ToLower(key)
		table, err := util.ValidChecker(key)

		if err != nil {
			requests.JSON(http.StatusOK, gin.H{
				"status":  http.StatusBadRequest,
				"error":   err.Error(),
				"message": "request key not valid",
			})
		} else {
			if table == "temporary" {
				paste := model.Temporary{Key: key}
				if err := paste.Get(); err != nil {
					if err.Error() == "record not found" {
						logger.Info(util.LoggerInfo(IP, "Access empty key: "+key))
						requests.JSON(http.StatusOK, gin.H{
							"status":  http.StatusNotFound,
							"error":   err.Error(),
							"message": fmt.Sprintf("key: %s not found", paste.Key),
						})
					} else {
						logger.Info(util.LoggerInfo(IP, "Query from db failed: "+err.Error()))
						requests.JSON(http.StatusInternalServerError, gin.H{
							"status":  http.StatusInternalServerError,
							"error":   err.Error(),
							"message": "query from db failed",
						})
					}
				} else {
					if paste.Password == "" || paste.Password == convert.String2md5(password) {
						logger.Info(util.LoggerInfo(IP, "Password accept"))
						if err := paste.Delete(); err != nil {
							requests.JSON(http.StatusInternalServerError, gin.H{
								"status":  http.StatusInternalServerError,
								"error":   err.Error(),
								"message": fmt.Sprintf("key: %s delete failed", paste.Key),
							})
						} else {
							jsonRequest := requests.DefaultQuery("json", "false")
							if jsonRequest == "false" { // API request
								logger.Info(util.LoggerInfo(IP, "jsonRequest: false"))
								requests.String(http.StatusOK, paste.Content)
							} else { // json request
								logger.Info(util.LoggerInfo(IP, "jsonRequest: true"))
								requests.JSON(http.StatusOK, gin.H{
									"status":  http.StatusOK,
									"lang":    paste.Lang,
									"content": paste.Content,
								})
							}
						}

					} else {
						logger.Info(util.LoggerInfo(IP, "Password wrong"))
						requests.JSON(http.StatusOK, gin.H{
							"status":  http.StatusUnauthorized,
							"error":   "wrong password",
							"message": "wrong password",
						})
					}
				}
			} else { // permanent
				paste := model.Permanent{Key: convert.String2uint(key)}
				if err := paste.Get(); err != nil {
					if err.Error() == "record not found" {
						logger.Info(util.LoggerInfo(IP, "Access empty key: "+key))
						requests.JSON(http.StatusOK, gin.H{
							"status":  http.StatusNotFound,
							"error":   err.Error(),
							"message": fmt.Sprintf("key: %d not found", paste.Key),
						})
					} else {
						logger.Info(util.LoggerInfo(IP, "Query from db failed: "+err.Error()))
						requests.JSON(http.StatusInternalServerError, gin.H{
							"status":  http.StatusInternalServerError,
							"error":   err.Error(),
							"message": "query from db failed",
						})
					}
				} else {
					if paste.Password == "" || paste.Password == convert.String2md5(password) {
						logger.Info(util.LoggerInfo(IP, "Password accept"))
						jsonRequest := requests.DefaultQuery("json", "false")
						if jsonRequest == "false" {
							logger.Info(util.LoggerInfo(IP, "jsonRequest: false"))
							requests.String(http.StatusOK, paste.Content)
						} else {
							logger.Info(util.LoggerInfo(IP, "jsonRequest: true"))
							requests.JSON(http.StatusOK, gin.H{
								"status":  http.StatusOK,
								"lang":    paste.Lang,
								"content": paste.Content,
							})
						}
					} else {
						logger.Info(util.LoggerInfo(IP, "Password wrong"))
						requests.JSON(http.StatusOK, gin.H{
							"status":  http.StatusUnauthorized,
							"error":   "wrong password",
							"message": "wrong password",
						})
					}
				}
			}
		}
	}
}

func notFound(requests *gin.Context) {
	requests.JSON(http.StatusNotFound, gin.H{
		"status":  http.StatusNotFound,
		"error":   "not found",
		"message": "no router founded",
	})
}
