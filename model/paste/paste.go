package paste

import (
	"github.com/PasteUs/PasteMeGoBackend/util"
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
	IPaste
	Key       string    `json:"key" gorm:"type:varchar(16);primary_key"` // 主键:索引
	Lang      string    `json:"lang" gorm:"type:varchar(16)"`            // 语言类型
	Content   string    `json:"content" gorm:"type:mediumtext"`          // 内容，最大长度为 16777215(2^24-1) 个字符
	Password  string    `json:"password" gorm:"type:varchar(32)"`        // 密码
	ClientIP  string    `json:"client_ip" gorm:"type:varchar(64)"`       // 用户 IP
	Username  string    `json:"username" gorm:"type:varchar(16)"`        // 用户名
	CreatedAt time.Time // 存储记录的创建时间
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

func (paste *AbstractPaste) checkPassword(password string) error {
	if paste.Password == "" || paste.Password == util.String2md5(password) {
		return nil
	}
	return ErrWrongPassword
}
