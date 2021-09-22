package paste

import (
	"github.com/PasteUs/PasteMeGoBackend/handler/session"
	"github.com/PasteUs/PasteMeGoBackend/logging"
	model "github.com/PasteUs/PasteMeGoBackend/model/paste"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func Create(context *gin.Context) {
	username := context.GetString(session.IdentityKey)

	body := requestBody{
		AbstractPaste: &model.AbstractPaste{
			ClientIP: context.ClientIP(),
			Username: username,
		},
	}

	if err := context.ShouldBindJSON(&body); err != nil {
		logging.Warn("bind body failed", context, zap.String("err", err.Error()))
		context.JSON(http.StatusOK, gin.H{
			"status":  http.StatusBadRequest,
			"message": "wrong param type",
		})
		return
	}

	if err := validator(body); err != nil {
		logging.Info("param validate failed", zap.String("err", err.Error()))
		context.JSON(http.StatusOK, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	if err := authenticator(body); err != nil {
		logging.Info("unauthorized request")
		context.JSON(http.StatusOK, gin.H{
			"status":  http.StatusUnauthorized,
			"message": err.Error(),
		})
		return
	}

	var paste model.IPaste

	if body.SelfDestruct {
		paste = &model.Temporary{
			AbstractPaste: body.AbstractPaste,
			ExpireMinute:  body.ExpireMinute,
			ExpireCount:   body.ExpireCount,
		}
	} else {
		paste = &model.Permanent{AbstractPaste: body.AbstractPaste}
	}

	if err := paste.Save(); err != nil {
		logging.Warn("save failed", context, zap.String("err", err.Error()))
		context.JSON(http.StatusOK, gin.H{
			"status":  http.StatusInternalServerError,
			"message": ErrSaveFailed.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"key":    paste.GetKey(),
	})
}

func Get(context *gin.Context) {
	key := strings.ToLower(context.Param("key"))

	var paste model.IPaste

	if err := keyValidator(key); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	abstractPaste := model.AbstractPaste{Key: key}

	if []rune(key)[0] == '0' {
		paste = &model.Temporary{AbstractPaste: &abstractPaste}
	} else {
		paste = &model.Permanent{AbstractPaste: &abstractPaste}
	}

	if err := paste.Get(context.DefaultQuery("password", "")); err != nil {
		var (
			status  int
			message string
		)

		switch err {
		case gorm.ErrRecordNotFound:
			status = http.StatusNotFound
			message = err.Error()
		case model.ErrWrongPassword:
			status = http.StatusForbidden
			message = err.Error()
		default:
			logging.Error("query from db failed", context, zap.String("err", err.Error()))
			status = http.StatusInternalServerError
			message = ErrQueryDBFailed.Error()
		}

		context.JSON(http.StatusOK, gin.H{
			"status":  status,
			"message": message,
		})

		return
	}

	if strings.Contains(context.GetHeader("Accept"), "json") {
		context.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"lang":    paste.GetLang(),
			"content": paste.GetContent(),
		})
	} else {
		context.String(http.StatusOK, paste.GetContent())
	}
}
