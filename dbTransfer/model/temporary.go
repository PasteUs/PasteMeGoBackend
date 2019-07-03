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
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

type Temporary struct {
	Key       string `json:"key" gorm:"type:varchar(16);primary_key;index:idx"`
	Lang      string `json:"lang" gorm:"type:varchar(16)"`
	Content   string `json:"content" gorm:"type:mediumtext"`
	Password  string `json:"password" gorm:"type:varchar(32)"`
	ClientIP  string `gorm:"type:varchar(64)"`
	CreatedAt time.Time
}

func (paste *Temporary) Save(db *gorm.DB) error {
	if paste.Content == "" {
		return errors.New("empty content")
	}
	if paste.Lang == "" {
		return errors.New("empty lang")
	}
	if strings.Contains(paste.Content, "#include") && paste.Lang == "plain" {
		paste.Lang = "cpp"
	}
	return db.Create(&paste).Error
}
