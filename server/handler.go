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
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func setPermanent(requests *gin.Context) {
	paste := model.Permanent{}
	if err := requests.Bind(&paste); err != nil {
		panic(err) // TODO
	} else {
		if err := paste.Save(); err != nil {
			requests.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		} else {
			requests.JSON(http.StatusCreated, gin.H{
				"status": http.StatusCreated,
				"Key":    paste.Key,
			})
		}
	}
}

func setTemporary(requests *gin.Context) {
	key := requests.Param("key")
	if key == "read_once" {
		paste := model.Temporary{Key: util.Generator()}
		if err := requests.Bind(&paste); err != nil {
			requests.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		} else {
			if err := paste.Save(); err != nil {
				requests.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
			} else {
				requests.JSON(http.StatusCreated, gin.H{
					"status": http.StatusCreated,
					"Key":    paste.Key,
				})
			}
		}
	} else {
		key = strings.ToLower(key)
		table, err := util.ValidChecker(key)
		if err != nil {
			if err.Error() == "wrong length" {
				requests.JSON(http.StatusBadRequest, gin.H{
					"message": "key's length should at least 3 and at most 8",
				})
			} else {
				requests.JSON(http.StatusBadRequest, gin.H{
					"message": "key should only contains digital or lowercase letters",
				})
			}
		} else {
			if table != "temporary" {
				requests.JSON(http.StatusBadRequest, gin.H{
					"message": "only temporary key can be specified",
				})
			} else {
				paste := model.Temporary{Key: key}
				if err := requests.Bind(&paste); err != nil {
					requests.JSON(http.StatusInternalServerError, gin.H{
						"message": err.Error(),
					})
				} else {
					if err := paste.Save(); err != nil {
						requests.JSON(http.StatusInternalServerError, gin.H{
							"message": err.Error(),
						})
					} else {
						requests.JSON(http.StatusCreated, gin.H{
							"status": http.StatusCreated,
							"Key":    paste.Key,
						})
					}
				}
			}
		}
	}
}

func get(requests *gin.Context) {
	token := requests.Param("token")
	if token == "" { // empty request
		requests.JSON(http.StatusOK, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Wrong params",
		})
	} else {
		key, password := util.Parse(token)
		key = strings.ToLower(key)
		table, err := util.ValidChecker(key)

		if err != nil {
			requests.JSON(http.StatusOK, gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
		} else {
			if table == "temporary" {
				paste := model.Temporary{Key: key}
				if err := paste.Get(); err != nil {
					requests.JSON(http.StatusOK, gin.H{
						"status":  http.StatusNotFound,
						"message": fmt.Sprintf("key: %s not found", paste.Key),
					})
				} else {
					if paste.Password == password {
						browser := requests.DefaultQuery("browser", "false")
						if browser == "false" { // API request
							requests.String(http.StatusOK, paste.Content)
						} else { // browser request
							requests.JSON(http.StatusOK, gin.H{
								"status":  http.StatusOK,
								"Lang":    paste.Lang,
								"Content": paste.Content,
							})
						}
					} else {
						requests.JSON(http.StatusOK, gin.H{
							"status":  http.StatusUnauthorized,
							"message": "Wrong password",
						})
					}
				}
			} else { // permanent
				paste := model.Permanent{Key: util.String2uint(key)}
				if err := paste.Get(); err != nil {
					if err.Error() == "record not found" {
						requests.JSON(http.StatusOK, gin.H{
							"status":  http.StatusNotFound,
							"message": fmt.Sprintf("key: %d not found", paste.Key),
						})
					} else {
						requests.JSON(http.StatusInternalServerError, gin.H{
							"message": err.Error(),
						})
					}
				} else {
					if paste.Password == password {
						browser := requests.DefaultQuery("browser", "empty")
						if browser == "empty" {
							requests.String(http.StatusOK, paste.Content)
						} else {
							requests.JSON(http.StatusOK, gin.H{
								"status":  http.StatusOK,
								"Lang":    paste.Lang,
								"Content": paste.Content,
							})
						}
					} else {
						requests.JSON(http.StatusOK, gin.H{
							"status":  http.StatusUnauthorized,
							"message": "Wrong password",
						})
					}
				}
			}
		}
	}
}
