/*
@File: model.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-11-20 16:27  bilibili     1.0.0       rewrite
*/
package model

import "time"

type PermanentUser struct{
	Key       uint64 `gorm:"primary_key"` // 主键:索引
	User
	// 存储记录的删除时间
	// 删除具有 DeletedAt 字段的记录，它不会从数据库中删除，但只将字段 DeletedAt 设置为当前时间，并在查询时无法找到记录
	DeletedAt *time.Time
}

// 成员函数，删除
func (paste *PermanentUser) Delete() error {
	return db.Delete(&paste, "`key` = ?", paste.Key).Error
}

// 成员函数，访问
func (paste *PermanentUser) Get() error {
	return db.Find(&paste, "`key` = ?", paste.Key).Error
}