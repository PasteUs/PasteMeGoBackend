package server

import (
    "fmt"
    "github.com/PasteUs/PasteMeGoBackend/model"
    "github.com/PasteUs/PasteMeGoBackend/util"
    "github.com/PasteUs/PasteMeGoBackend/util/convert"
    "github.com/PasteUs/PasteMeGoBackend/util/generator"
    "github.com/gin-gonic/gin"
    "github.com/wonderivan/logger"
    "net/http"
    "strings"
)

// 创建一个永久的 Paste, key 是自增键
func permanentCreator(requests *gin.Context) {
    IP := requests.ClientIP() // 用户 IP
    paste := model.Permanent{
        AbstractPaste: &model.AbstractPaste{
            ClientIP: IP,
        },
    }
    // 绑定请求参数
    if err := requests.ShouldBindJSON(&paste); err != nil {
        logger.Error(util.LogFormat(IP, "Bind failed: %s", err.Error()))
        requests.JSON(http.StatusInternalServerError, gin.H{
            "status":  http.StatusInternalServerError,
            "message": "bind failed",
            "error":   err.Error(),
        })
    } else {
        if err := paste.Save(); err != nil {
            logger.Error(util.LogFormat(IP, "Save failed: %s", err.Error()))
            requests.JSON(http.StatusInternalServerError, gin.H{
                "status":  http.StatusInternalServerError,
                "message": "save failed",
                "error":   err.Error(),
            })
        } else {
            logger.Info(util.LogFormat(IP, "Create an permanent paste with key: "+convert.Uint2string(paste.Key)))

            requests.JSON(http.StatusCreated, gin.H{
                "status": http.StatusCreated,
                "key":    paste.Key,
            })
        }
    }
}

// 创建一个阅后即焚的 Paste, key 是指定的
func temporaryCreator(requests *gin.Context) {
    IP, key := requests.ClientIP(), requests.Param("key")
    key = strings.ToLower(key) // 进行大写到小写的转换
    table, err := util.ValidChecker(key)
    if err != nil {
        if err.Error() == "wrong length" {
            logger.Warn(util.LogFormat(IP, "Trying to create temporary paste with key: %s", key))
            requests.JSON(http.StatusOK, gin.H{
                "status":  http.StatusBadRequest,
                "error":   err.Error(),
                "message": "key's length should at least 3 and at most 8",
            })
        } else {
            logger.Warn(util.LogFormat(IP, "Trying to create temporary paste with key: %s", key))
            requests.JSON(http.StatusOK, gin.H{
                "status":  http.StatusBadRequest,
                "error":   err.Error(),
                "message": "temporary key should only contains digits and lowercase letters, at least one alpha is required",
            })
        }
    } else {
        if table != "temporary" {
            logger.Warn(util.LogFormat(IP, "Trying to create temporary paste with key: %s", key))
            requests.JSON(http.StatusOK, gin.H{
                "status":  http.StatusBadRequest,
                "error":   "wrong key type",
                "message": "temporary key should only contains digits and lowercase letters, at least one alpha is required",
            })
        } else {
            paste := model.Temporary{
                Key: key,
                AbstractPaste: &model.AbstractPaste{
                    ClientIP: requests.ClientIP(),
                },
            }
            if err := requests.ShouldBindJSON(&paste); err != nil {
                logger.Error(util.LogFormat(IP, "Bind failed: %s", err.Error()))
                requests.JSON(http.StatusInternalServerError, gin.H{
                    "status":  http.StatusInternalServerError,
                    "error":   err.Error(),
                    "message": "bind failed",
                })
            } else {
                if err := paste.Save(); err != nil {
                    logger.Error(util.LogFormat(IP, "Save failed: %s", err.Error()))
                    requests.JSON(http.StatusInternalServerError, gin.H{
                        "status":  http.StatusInternalServerError,
                        "error":   err.Error(),
                        "message": "save failed",
                    })
                } else {
                    logger.Info(util.LogFormat(IP, "Create an temporary paste with key: "+paste.Key))
                    requests.JSON(http.StatusCreated, gin.H{
                        "status": http.StatusCreated,
                        "key":    paste.Key,
                    })
                }
            }
        }
    }
}

// 创建一个阅后即焚的 Paste, key 是随机的
func readOnceCreator(requests *gin.Context) {
    IP := requests.ClientIP()
    paste := model.Temporary{
        Key: generator.Generator(),
        AbstractPaste: &model.AbstractPaste{
            ClientIP: IP,
        },
    }
    if err := requests.ShouldBindJSON(&paste); err != nil {
        logger.Error(util.LogFormat(IP, "Bind failed: %s", err.Error()))
        requests.JSON(http.StatusInternalServerError, gin.H{
            "status":  http.StatusInternalServerError,
            "error":   err.Error(),
            "message": "bind failed",
        })
    } else {
        if err := paste.Save(); err != nil {
            logger.Error(util.LogFormat(IP, "Save failed: %s", err.Error()))
            requests.JSON(http.StatusInternalServerError, gin.H{
                "status":  http.StatusInternalServerError,
                "error":   err.Error(),
                "message": "save failed",
            })
        } else {
            logger.Info(util.LogFormat(IP, "Create an once paste with key: "+paste.Key))
            requests.JSON(http.StatusCreated, gin.H{
                "status": http.StatusCreated,
                "key":    paste.Key,
            })
        }
    }
}

// 访问未加密的 Paste, token 为 <Paste ID>
// 访问加密的 Paste, token 为 <Paste ID>,<Password>
func query(requests *gin.Context) {
    IP, token := requests.ClientIP(), requests.Param("token")
    if token == "" { // 空的 token
        requests.JSON(http.StatusOK, gin.H{
            "status":  http.StatusBadRequest,
            "error":   "empty token",
            "message": "wrong params",
        })
    } else {
        key, password := util.Parse(token)   // 分离出 key 和 password
        key = strings.ToLower(key)           // 进行大写到小写的转换
        table, err := util.ValidChecker(key) // 正则匹配

        if err != nil {
            requests.JSON(http.StatusOK, gin.H{
                "status":  http.StatusBadRequest,
                "error":   err.Error(),
                "message": "request key not valid",
            })
        } else {
            var paste model.IPaste
            if table == "temporary" {
                paste = &model.Temporary{Key: key}
            } else {
                paste = &model.Permanent{Key: convert.String2uint(key)}
            }

            if err := paste.Get(); err != nil {
                if err.Error() == "record not found" {
                    logger.Info(util.LogFormat(IP, "Access empty key: %s", key))
                    requests.JSON(http.StatusOK, gin.H{
                        "status":  http.StatusNotFound,
                        "error":   err.Error(),
                        "message": fmt.Sprintf("key: %s not found", key),
                    })
                } else {
                    logger.Info(util.LogFormat(IP, "Query from db failed: %s", err.Error()))
                    requests.JSON(http.StatusInternalServerError, gin.H{
                        "status":  http.StatusInternalServerError,
                        "error":   err.Error(),
                        "message": "query from db failed",
                    })
                }
            } else {
                if paste.GetPassword() == "" || paste.GetPassword() == convert.String2md5(password) { // 密码为空或者密码正确
                    logger.Info(util.LogFormat(IP, "Password accept"))
                    if table == "temporary" {
                        if err := paste.Delete(); err != nil {
                            requests.JSON(http.StatusInternalServerError, gin.H{
                                "status":  http.StatusInternalServerError,
                                "error":   err.Error(),
                                "message": fmt.Sprintf("key: %s delete failed", key),
                            })
                            return
                        }
                    }

                    jsonRequest := requests.DefaultQuery("json", "false")
                    if jsonRequest == "false" { // raw request
                        logger.Info(util.LogFormat(IP, "jsonRequest: false"))
                        requests.String(http.StatusOK, paste.GetContent())
                    } else { // json request
                        logger.Info(util.LogFormat(IP, "jsonRequest: true"))
                        requests.JSON(http.StatusOK, gin.H{
                            "status":  http.StatusOK,
                            "lang":    paste.GetLang(),
                            "content": paste.GetContent(),
                        })
                    }
                } else {
                    logger.Info(util.LogFormat(IP, "Password wrong")) // 密码错误
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

func notFoundHandler(requests *gin.Context) {
    requests.JSON(http.StatusNotFound, gin.H{
        "status":  http.StatusNotFound,
        "error":   "not found",
        "message": "no router founded",
    })
}

func beat(requests *gin.Context) {
    method := requests.DefaultQuery("method", "none")
    if method == "beat" {
        requests.JSON(http.StatusOK, gin.H{
            "status": http.StatusOK,
        })
    } else {
        notFoundHandler(requests)
    }
}
