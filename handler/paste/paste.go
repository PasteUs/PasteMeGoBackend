package paste

import (
	"github.com/PasteUs/PasteMeGoBackend/handler/common"
	"github.com/PasteUs/PasteMeGoBackend/handler/session"
	"github.com/PasteUs/PasteMeGoBackend/logging"
	model "github.com/PasteUs/PasteMeGoBackend/model/paste"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// Create 创建一贴
// @Summary 创建永久存储或者是自我销毁的一贴
// @Description 只有在登陆的状态下才能创建永久的一贴
// @Tags Paste
// @Accept json
// @Produce json
// @Param Authorization header string false "登陆的 Token"
// @Param data body CreateRequest true "请求数据"
// @Success 201 {object} CreateResponse
// @Failure default {object} common.ErrorResponse
// @Router /paste/ [post]
func Create(context *gin.Context) {
	username := context.GetString(session.IdentityKey)

	body := CreateRequest{
		AbstractPaste: &model.AbstractPaste{
			ClientIP: context.ClientIP(),
			Username: username,
		},
	}

	if err := context.ShouldBindJSON(&body); err != nil {
		logging.Warn("bind body failed", context, zap.String("err", err.Error()))
		common.Error(context, http.StatusBadRequest, ErrWrongParamType)
		return
	}

	if err := validator(body); err != nil {
		logging.Info("param validate failed", zap.String("err", err.Error()))
		common.Error(context, http.StatusBadRequest, err)
		return
	}

	if err := authenticator(body); err != nil {
		logging.Info("unauthorized request")
		common.Error(context, http.StatusUnauthorized, err)
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
		common.Error(context, http.StatusInternalServerError, err)
		return
	}

	context.JSON(http.StatusCreated, CreateResponse{
		Response: &common.Response{Status: http.StatusCreated},
		Key:      paste.GetKey(),
	})
}

// Get godoc
// @Summary 读取一贴
// @Description 如果不指定 Accept: application/json 的话，默认会返回 text/plain 格式的 content
// @Tags Paste
// @Accept json
// @Produce json
// @Param Accept header string false "响应格式" default("text/plain")
// @Param key path string true "索引"
// @Success 201 {object} GetResponse
// @Failure default {object} common.ErrorResponse
// @Router /paste/{key} [get]
func Get(context *gin.Context) {
	key := strings.ToLower(context.Param("key"))

	var paste model.IPaste

	if err := keyValidator(key); err != nil {
		common.Error(context, http.StatusBadRequest, err)
		return
	}

	abstractPaste := model.AbstractPaste{Key: key}

	if []rune(key)[0] == '0' {
		paste = &model.Temporary{AbstractPaste: &abstractPaste}
	} else {
		paste = &model.Permanent{AbstractPaste: &abstractPaste}
	}

	if err := paste.Get(context.DefaultQuery("password", "")); err != nil {
		var status int
		switch err {
		case gorm.ErrRecordNotFound:
			status = http.StatusNotFound
		case model.ErrWrongPassword:
			status = http.StatusForbidden
		default:
			logging.Error("query from db failed", context, zap.String("err", err.Error()))
			status = http.StatusInternalServerError
		}

		common.Error(context, status, err)
		return
	}

	if strings.Contains(context.GetHeader("Accept"), "json") {
		context.JSON(http.StatusOK, GetResponse{
			Response: &common.Response{
				Status: http.StatusOK,
			},
			Lang: paste.GetLang(),
			Content: paste.GetContent(),
		})
	} else {
		context.String(http.StatusOK, paste.GetContent())
	}
}
