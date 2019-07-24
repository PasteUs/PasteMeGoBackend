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
	"github.com/LucienShui/PasteMeBackend/util/convert"
	"time"
)

type Permanent struct {
	Key       uint64 `gorm:"primary_key;index:idx"`
	Lang      string `json:"lang" gorm:"type:varchar(16)"`
	Content   string `json:"content" gorm:"type:mediumtext"`
	Password  string `json:"password" gorm:"type:varchar(32)"`
	ClientIP  string `gorm:"type:varchar(64)"`
	CreatedAt time.Time
	DeletedAt *time.Time
}

func (paste *Permanent) Save() error {
	if paste.Content == "" {
		return errors.New("empty content")
	}
	if paste.Lang == "" {
		return errors.New("empty lang")
	}
	if paste.Password != "" {
		paste.Password = convert.String2md5(paste.Password)
	}
	return db.Create(&paste).Error
}

func (paste *Permanent) Delete() error {
	return db.Delete(&paste, "`key` = ?", paste.Key).Error
}

func (paste *Permanent) Get() error {
	return db.Find(&paste, "`key` = ?", paste.Key).Error
}
