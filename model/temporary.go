/*
@File: paste.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-06-23 14:03  Lucien     1.0         None
*/
package model

import (
	"errors"
	"github.com/PasteUs/PasteMeGoBackend/util/convert"
	"time"
)

// 临时
type Temporary struct {
	Key       string `json:"key" gorm:"type:varchar(16);primary_key;index:idx"` // 主键:索引
	Lang      string `json:"lang" gorm:"type:varchar(16)"`  // 语言类型
	Content   string `json:"content" gorm:"type:mediumtext"` // 内容，最大长度为 16777215(2^24-1) 个字符
	Password  string `json:"password" gorm:"type:varchar(32)"` // 密码
	ClientIP  string `gorm:"type:varchar(64)"` // 用户 IP
	CreatedAt time.Time // 存储记录的创建时间
}

// 成员函数，保存
func (paste *Temporary) Save() error {
	if paste.Content == "" {
		return errors.New("empty content") //内容为空，返回错误信息 "empty content"
	}
	if paste.Lang == "" {
		return errors.New("empty lang") // 语言类型为空，返回错误信息 "empty lang"
	}
	if paste.Password != "" {
		paste.Password = convert.String2md5(paste.Password)
	}
	return db.Create(&paste).Error
}

// 成员函数，删除
func (paste *Temporary) Delete() error {
	return db.Delete(&paste, "`key` = ?", paste.Key).Error
}

// 成员函数，查看
func (paste *Temporary) Get() error {
	return db.Find(&paste, "`key` = ?", paste.Key).Error
}

func Exist(key string) bool {
	count := uint8(0)
	db.Model(&Temporary{}).Where("`key` = ?", key).Count(&count)
	return count > 0
}
