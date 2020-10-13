/*
@File: paste.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-06-23 14:03  Lucien     1.0         None
2020-10-12 20:06  Lx200916   1.1         Add Burn after reading
*/
package model

import (
	"errors"
	"github.com/PasteUs/PasteMeGoBackend/config"
	"github.com/PasteUs/PasteMeGoBackend/util/convert"
	"time"
)

// 永久
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

	// 过期时间,为Time结构体,标志的是消息删除的日期时间,为零值时不过期.利用自定义的类型进行自定义解析转换.
	// 过期后查询时无法找到记录,在某个指定时间由定时任务 Clean() 统一回收
	DeadLine *DeadLine `gorm:"type:datetime"`
}

//成员函数，创建
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

// 成员函数，删除
func (paste *Permanent) Delete() error {
	return db.Delete(&paste, "`key` = ?", paste.Key).Error
}

// 成员函数，访问
func (paste *Permanent) Get() error {
	if config.Get().Database.Type == "mysql" {
		return db.Find(&paste, "`key` = ? and (dead_line > (select now()) or dead_line is NULL)", paste.Key).Error

	} else {
		return db.Find(&paste, "`key` = ? and (dead_line > DATETIME('now') or dead_line is NULL)", paste.Key).Error
	}
}

//成员函数,自动删除
func Clean() error {
	if config.Get().Database.Type == "mysql" {
		return db.Exec("DELETE FROM Permanents WHERE dead_line < (SELECT now())").Error
	} else {
		return db.Exec("DELETE FROM Permanents WHERE dead_line < DATETIME('now')").Error
	}

}
