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
	"net/http"
	"strings"
)

func setPermanent(requests *gin.Context) {
	paste := model.Permanent{
		ClientIP: requests.ClientIP(),
	}
	if err := requests.ShouldBindJSON(&paste); err != nil {
		requests.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "bind failed",
			"error":   err.Error(),
		})
	} else {
		if err := paste.Save(); err != nil {
			requests.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "save failed",
				"error":   err.Error(),
			})
		} else {
			requests.JSON(http.StatusCreated, gin.H{
				"status": http.StatusCreated,
				"key":    paste.Key,
			})
		}
	}
}

func setTemporary(requests *gin.Context) {
	key := requests.Param("key")
	if key == "read_once" {
		paste := model.Temporary{
			Key:      util.Generator(),
			ClientIP: requests.ClientIP(),
		}
		if err := requests.ShouldBindJSON(&paste); err != nil {
			requests.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"error":   err.Error(),
				"message": "bind failed",
			})
		} else {
			if err := paste.Save(); err != nil {
				requests.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"error":   err.Error(),
					"message": "save failed",
				})
			} else {
				requests.JSON(http.StatusCreated, gin.H{
					"status": http.StatusCreated,
					"key":    paste.Key,
				})
			}
		}
	} else {
		key = strings.ToLower(key)
		table, err := util.ValidChecker(key)
		if err != nil {
			if err.Error() == "wrong length" {
				requests.JSON(http.StatusOK, gin.H{
					"status":  http.StatusBadRequest,
					"error":   err.Error(),
					"message": "key's length should at least 3 and at most 8",
				})
			} else {
				requests.JSON(http.StatusOK, gin.H{
					"status":  http.StatusBadRequest,
					"error":   err.Error(),
					"message": "key should only contains digital or lowercase letters",
				})
			}
		} else {
			if table != "temporary" {
				requests.JSON(http.StatusOK, gin.H{
					"status":  http.StatusBadRequest,
					"error":   "wrong key type",
					"message": "only temporary key can be specified",
				})
			} else {
				paste := model.Temporary{
					Key:      key,
					ClientIP: requests.ClientIP(),
				}
				if err := requests.ShouldBindJSON(&paste); err != nil {
					requests.JSON(http.StatusInternalServerError, gin.H{
						"status":  http.StatusInternalServerError,
						"error":   err.Error(),
						"message": "bind failed",
					})
				} else {
					if err := paste.Save(); err != nil {
						requests.JSON(http.StatusInternalServerError, gin.H{
							"status":  http.StatusInternalServerError,
							"error":   err.Error(),
							"message": "save failed",
						})
					} else {
						requests.JSON(http.StatusCreated, gin.H{
							"status": http.StatusCreated,
							"key":    paste.Key,
						})
					}
				}
			}
		}
	}
}

func get(requests *gin.Context) {
	token := requests.Param("token")
	if token == "" { // empty token
		requests.JSON(http.StatusOK, gin.H{
			"status":  http.StatusBadRequest,
			"error":   "empty token",
			"message": "Wrong params",
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
					requests.JSON(http.StatusOK, gin.H{
						"status":  http.StatusNotFound,
						"error":   err.Error(),
						"message": fmt.Sprintf("key: %s not found", paste.Key),
					})
				} else {
					if paste.Password == "" || paste.Password == convert.String2md5(password) {
						if err := paste.Delete(); err != nil {
							requests.JSON(http.StatusInternalServerError, gin.H{
								"status":  http.StatusInternalServerError,
								"error":   err.Error(),
								"message": fmt.Sprintf("key: %s delete failed", paste.Key),
							})
						} else {
							browser := requests.DefaultQuery("browser", "false")
							if browser == "false" { // API request
								requests.String(http.StatusOK, paste.Content)
							} else { // browser request
								requests.JSON(http.StatusOK, gin.H{
									"status":  http.StatusOK,
									"lang":    paste.Lang,
									"content": paste.Content,
								})
							}
						}

					} else {
						requests.JSON(http.StatusOK, gin.H{
							"status":  http.StatusUnauthorized,
							"error":   "wrong password",
							"message": "Wrong password",
						})
					}
				}
			} else { // permanent
				paste := model.Permanent{Key: convert.String2uint(key)}
				if err := paste.Get(); err != nil {
					if err.Error() == "record not found" {
						requests.JSON(http.StatusOK, gin.H{
							"status":  http.StatusNotFound,
							"error":   err.Error(),
							"message": fmt.Sprintf("key: %d not found", paste.Key),
						})
					} else {
						requests.JSON(http.StatusInternalServerError, gin.H{
							"status":  http.StatusInternalServerError,
							"error":   err.Error(),
							"message": "query from db failed",
						})
					}
				} else {
					if paste.Password == "" || paste.Password == convert.String2md5(password) {
						browser := requests.DefaultQuery("browser", "empty")
						if browser == "empty" {
							requests.String(http.StatusOK, paste.Content)
						} else {
							requests.JSON(http.StatusOK, gin.H{
								"status":  http.StatusOK,
								"lang":    paste.Lang,
								"content": paste.Content,
							})
						}
					} else {
						requests.JSON(http.StatusOK, gin.H{
							"status":  http.StatusUnauthorized,
							"error":   "wrong password",
							"message": "Wrong password",
						})
					}
				}
			}
		}
	}
}
