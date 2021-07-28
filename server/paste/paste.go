package paste

import (
    model "github.com/PasteUs/PasteMeGoBackend/model/paste"
    "github.com/PasteUs/PasteMeGoBackend/util"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "net/http"
)


func validator(expireType string, expiration uint64) error {
    if expireType == "" || expiration == 0 {
        return ErrEmptyExpireTypeOrExpiration
    }

    if expiration <= 0 {
        return ErrZeroExpiration
    }

    if expireType == model.EnumTime {
        if expiration > model.OneMonth {
            return ErrExpirationGreaterThanMonth
        }
    } else if expireType == model.EnumCount {
        if expiration > model.MaxCount {
            return ErrExpirationGreaterThanMaxCount
        }
    } else {
        return ErrInvalidExpireType
    }
    return nil
}

func Create(context *gin.Context) {
    namespace := context.Param("namespace")
    util.Info("create paste", context, zap.String("namespace", namespace))

    body := struct {
        *model.AbstractPaste
        SelfDestruct bool   `json:"self_destruct"`
        ExpireType   string `json:"expire_type"`
        Expiration   uint64 `json:"expiration"`
    }{
        AbstractPaste: &model.AbstractPaste{
            ClientIP: context.ClientIP(),
        },
    }

    if err := context.ShouldBindJSON(&body); err != nil {
        util.Warn("bind body failed", context, zap.String("err", err.Error()))
        context.JSON(http.StatusBadRequest, gin.H{
            "status":  http.StatusBadRequest,
            "message": "wrong param type",
        })
        return
    }

    var paste model.IPaste

    if body.SelfDestruct {
        if err := validator(body.ExpireType, body.Expiration); err != nil {
            util.Info("temporary param validate failed", zap.String("err", err.Error()))
            context.JSON(http.StatusBadRequest, gin.H{
                "status":  http.StatusBadRequest,
                "message": err.Error(),
            })
        }

        paste = &model.Temporary{
            Key:           model.Generator(),
            Namespace:     namespace,
            AbstractPaste: body.AbstractPaste,
            ExpireType:    body.ExpireType,
            Expiration:    body.Expiration,
        }
    } else {
        paste = &model.Permanent{Namespace: namespace, AbstractPaste: body.AbstractPaste}
    }

    if err := paste.Save(); err != nil {
        util.Warn("save failed", context, zap.String("err", err.Error()))
        context.JSON(http.StatusInternalServerError, gin.H{
            "status":  http.StatusInternalServerError,
            "message": "save failed",
        })
        return
    }

    context.JSON(http.StatusCreated, gin.H{
        "status":    http.StatusCreated,
        "key":       paste.GetKey(),
        "namespace": paste.GetNamespace(),
    })
}

func Get(context *gin.Context) {
    namespace, key := context.Param("namespace"), context.Param("key")
    util.Info("test", context, zap.String("namespace", namespace), zap.String("key", key))
    context.JSON(http.StatusOK, gin.H{
        "status": http.StatusOK,
        "token":  "Hello World!",
    })
}
