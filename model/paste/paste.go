package paste

import (
	"crypto/md5"
	"fmt"
	"github.com/PasteUs/PasteMeGoBackend/handler/common"
	"github.com/PasteUs/PasteMeGoBackend/model/dao"
	"time"
)

type IPaste interface {
	Save() error
	Get(string) error
	Delete() error
	GetKey() string
	GetContent() string
	GetLang() string
}

type AbstractPaste struct {
	Key       string    `json:"key" swaggerignore:"true" gorm:"type:varchar(16);primaryKey"` // 主键:索引
	Lang      string    `json:"lang" example:"plain" gorm:"type:varchar(16)"`                // 语言类型
	Content   string    `json:"content" example:"Hello World!" gorm:"type:mediumtext"`       // 内容，最大长度为 16777215(2^24-1) 个字符
	Password  string    `json:"password" example:"" gorm:"type:varchar(32)"`                 // 密码
	ClientIP  string    `json:"client_ip" swaggerignore:"true" gorm:"type:varchar(64)"`      // 用户 IP
	Username  string    `json:"username" swaggerignore:"true" gorm:"type:varchar(16)"`       // 用户名
	CreatedAt time.Time `swaggerignore:"true"`                                               // 存储记录的创建时间
}

func (paste *AbstractPaste) GetKey() string {
	return paste.Key
}

func (paste *AbstractPaste) GetContent() string {
	return paste.Content
}

func (paste *AbstractPaste) GetLang() string {
	return paste.Lang
}

func hash(text string) string {
	if text == "" {
		return text
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(text)))
}

func (paste *AbstractPaste) checkPassword(password string) *common.ErrorResponse {
	if paste.Password == hash(password) {
		return nil
	}
	return common.ErrWrongPassword
}

func exist(key string, model interface{}) bool {
	count := int64(0)
	dao.DB.Model(model).Where("`key` = ?", key).Count(&count)
	return count > 0
}
