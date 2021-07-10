package model

import (
	"errors"
	"github.com/PasteUs/PasteMeGoBackend/util/convert"
	"time"
)

// Permanent 永久
type Permanent struct {
	Key       uint64    `gorm:"primary_key"`                      // 主键:索引
	Lang      string    `json:"lang" gorm:"type:varchar(16)"`     // 语言类型
	Content   string    `json:"content" gorm:"type:mediumtext"`   // 内容，最大长度为 16777215(2^24-1) 个字符
	Password  string    `json:"password" gorm:"type:varchar(32)"` // 密码
	ClientIP  string    `gorm:"type:varchar(64)"`                 // 用户 IP
	CreatedAt time.Time // 存储记录的创建时间
	// 存储记录的删除时间
	// 删除具有 DeletedAt 字段的记录，它不会从数据库中删除，但只将字段 DeletedAt 设置为当前时间，并在查询时无法找到记录
	DeletedAt *time.Time
}

// Save 成员函数，创建
func (paste *Permanent) Save() error {
	if paste.Content == "" {
		return errors.New("empty content") // 内容为空，返回错误信息 "empty content"
	}
	if paste.Lang == "" {
		return errors.New("empty lang") // 语言类型为空，返回错误信息 "empty lang"
	}
	if paste.Password != "" {
		paste.Password = convert.String2md5(paste.Password) // 加密存储，设置密码
	}
	return db.Create(&paste).Error
}

// Delete 成员函数，删除
func (paste *Permanent) Delete() error {
	return db.Delete(&paste, "`key` = ?", paste.Key).Error
}

// Get 成员函数，访问
func (paste *Permanent) Get() error {
	return db.Find(&paste, "`key` = ?", paste.Key).Error
}
