package paste

import (
	"github.com/PasteUs/PasteMeGoBackend/handler/common"
	"github.com/PasteUs/PasteMeGoBackend/handler/token"
	model "github.com/PasteUs/PasteMeGoBackend/model/paste"
	"regexp"
)

var (
	validLang  = []string{"plain", "cpp", "java", "python", "bash", "markdown", "json", "go"}
	keyPattern = regexp.MustCompile("^[0-9a-z]{8}$")
)

type CreateRequest struct {
	*model.AbstractPaste
	SelfDestruct bool   `json:"self_destruct" example:"true"` // 是否自我销毁
	ExpireSecond uint64 `json:"expire_second" example:"300"`    // 创建若干秒后自我销毁
	ExpireCount  uint64 `json:"expire_count" example:"1"`     // 访问若干次后自我销毁
}

type CreateResponse struct {
	*common.Response
	Key string `json:"key" example:"a1b2c3d4"`
}

type GetResponse struct {
	*common.Response
	Lang    string `json:"lang" example:"plain"`
	Content string `json:"content" example:"Hello World!"`
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func validator(body CreateRequest) *common.ErrorResponse {
	if body.Content == "" {
		return common.ErrEmptyContent // 内容为空，返回错误信息 "empty content"
	}
	if body.Lang == "" {
		return common.ErrEmptyLang // 语言类型为空，返回错误信息 "empty lang"
	}
	if !contains(validLang, body.Lang) {
		return common.ErrInvalidLang
	}

	if body.SelfDestruct {
		if body.ExpireSecond <= 0 {
			return common.ErrZeroExpireSecond
		}
		if body.ExpireCount <= 0 {
			return common.ErrZeroExpireCount
		}

		if body.ExpireSecond > model.OneMonth {
			return common.ErrExpireSecondGreaterThanMonth
		}
		if body.ExpireCount > model.MaxCount {
			return common.ErrExpireCountGreaterThanMaxCount
		}
	}
	return nil
}

func authenticator(body CreateRequest) *common.ErrorResponse {
	if body.Username == token.Nobody {
		if !body.SelfDestruct {
			return common.ErrUnauthorized
		}
		if body.ExpireCount > 1 {
			return common.ErrUnauthorized
		}
		if body.ExpireSecond > 5 * 60 {
			return common.ErrUnauthorized
		}
	}
	return nil
}

func keyValidator(key string) *common.ErrorResponse {
	if len(key) != 8 {
		return common.ErrInvalidKeyLength // key's length should at least 3 and at most 8
	}
	if flag := keyPattern.MatchString(key); !flag {
		return common.ErrInvalidKeyFormat
	}
	return nil
}
