/*
@File: handler.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Desciption
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
)

func set(requests *gin.Context) {
	temporary := model.Temporary{}
	if err := requests.Bind(&temporary); err != nil {
		panic(err) // TODO
	} else {
		if temporary.Key == "" { // permanent
			permanent := model.Permanent{}
			permanent.Load(temporary)
			if err := permanent.Save(); err != nil {
				panic(err) // TODO
			} else {
				requests.JSON(http.StatusCreated, gin.H{
					"status": http.StatusCreated,
					"Key":    permanent.Key,
				})
			}
		} else if temporary.Key == "read_once" {
			temporary.Key = util.Generator()
			if err := temporary.Save(); err != nil {
				panic(err) // TODO
			} else {
				requests.JSON(http.StatusCreated, gin.H{
					"status": http.StatusCreated,
					"Key":    temporary.Key,
				})
			}
		} else {
			table, err := util.ValidChecker(temporary.Key)
			if err != nil {
				panic(err) // TODO
			} else {
				if table != "temporary" {
					// TODO
				} else {
					if err := temporary.Save(); err != nil {
						panic(err) // TODO
					} else {
						requests.JSON(http.StatusCreated, gin.H{
							"status": http.StatusCreated,
							"Key":    temporary.Key,
						})
					}
				}
			}
		}
	}
}

func get(requests *gin.Context) {
	token := requests.DefaultQuery("token", "")
	if token == "" { // empty request
		requests.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"message": "Wrong params",
		})
	} else {
		key, password := util.Parse(token)
		table, err := util.ValidChecker(key)

		if err != nil {
			requests.JSON(http.StatusOK, gin.H{
				"status": http.StatusBadRequest,
				"message": err.Error(),
			})
		} else {
			if table == "temporary" {
				paste := model.Temporary{Key: key}
				if err := paste.Get(); err != nil {
					requests.JSON(http.StatusOK, gin.H{
						"status": http.StatusNotFound,
						"message": fmt.Sprintf("key: %s not found", paste.Key),
					})
				} else {
					if paste.Password == password {
						browser := requests.DefaultQuery("browser", "false")
						if browser == "false" { // API request
							requests.String(http.StatusOK, paste.Content)
						} else { // browser request
							requests.JSON(http.StatusOK, gin.H{
								"status": http.StatusOK,
								"Lang": paste.Lang,
								"Content": paste.Content,
							})
						}
					} else {
						requests.JSON(http.StatusOK, gin.H{
							"status": http.StatusUnauthorized,
							"message": "Wrong password",
						})
					}
				}
			} else { // permanent
				paste := model.Permanent{Key: util.String2uint(key)}
				if err := paste.Get(); err != nil {
					if err.Error() == "record not found" {
						requests.JSON(http.StatusOK, gin.H{
							"status": http.StatusNotFound,
							"message": fmt.Sprintf("key: %d not found", paste.Key),
						})
					} else {
						panic(err) // TODO
					}
				} else {
					if paste.Password == password {
						browser := requests.DefaultQuery("browser", "empty")
						if browser == "empty" {
							requests.String(http.StatusOK, paste.Content)
						} else {
							requests.JSON(http.StatusOK, gin.H{
								"status": http.StatusOK,
								"Lang": paste.Lang,
								"Content": paste.Content,
							})
						}
					} else {
						requests.JSON(http.StatusOK, gin.H{
							"status": http.StatusUnauthorized,
							"message": "Wrong password",
						})
					}
				}
			}
		}
	}
}
