package paste

import (
	"github.com/PasteUs/PasteMeGoBackend/handler/session"
	model "github.com/PasteUs/PasteMeGoBackend/model/paste"
	"regexp"
)

var (
	validLang  = []string{"plain", "cpp", "java", "python", "bash", "markdown", "json", "go"}
	keyPattern = regexp.MustCompile("^[0-9a-z]{8}$")
)

type requestBody struct {
	*model.AbstractPaste
	SelfDestruct bool   `json:"self_destruct"`
	ExpireMinute uint64 `json:"expire_minute"`
	ExpireCount  uint64 `json:"expire_count"`
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func validator(body requestBody) error {
	if body.Content == "" {
		return ErrEmptyContent // 内容为空，返回错误信息 "empty content"
	}
	if body.Lang == "" {
		return ErrEmptyLang // 语言类型为空，返回错误信息 "empty lang"
	}
	if !contains(validLang, body.Lang) {
		return ErrInvalidLang
	}

	if body.SelfDestruct {
		if body.ExpireMinute <= 0 {
			return ErrZeroExpireMinute
		}
		if body.ExpireCount <= 0 {
			return ErrZeroExpireCount
		}

		if body.ExpireMinute > model.OneMonth {
			return ErrExpireMinuteGreaterThanMonth
		}
		if body.ExpireCount > model.MaxCount {
			return ErrExpireCountGreaterThanMaxCount
		}
	}
	return nil
}

func authenticator(body requestBody) error {
	if body.Username == session.Nobody {
		if !body.SelfDestruct {
			return ErrUnauthorized
		}
		if body.ExpireCount > 1 {
			return ErrUnauthorized
		}
		if body.ExpireMinute > 5 {
			return ErrUnauthorized
		}
	}
	return nil
}

func keyValidator(key string) error {
	if len(key) != 8 {
		return ErrInvalidKeyLength // key's length should at least 3 and at most 8
	}
	if flag := keyPattern.MatchString(key); !flag {
		return ErrInvalidKeyFormat
	}
	return nil
}
