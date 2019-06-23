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
	"time"
)

type Temporary struct {
	Key       string `gorm:"type:varchar(17);primary_key;index:idx"`
	Lang      string `gorm:"type:varchar(17)"`
	Content   string `gorm:"type:mediumtext"`
	Password  string `gorm:"type:varchar(17)"`
	CreatedAt time.Time
}

func (paste *Temporary) Save() error {
	return db.Create(&paste).Error
}

func (paste *Temporary) Delete() error {
	return db.Delete(&paste, "`key` = ?", paste.Key).Error
}

func (paste *Temporary) Get() error {
	if err := db.Find(&paste, "`key` = ?", paste.Key).Error; err != nil {
		return err
	}
	return paste.Delete()
}

func Exist(key string) bool {
	count := uint8(0)
	db.Model(&Temporary{}).Where("`key` = ?", key).Count(&count)
	return count > 0
}
