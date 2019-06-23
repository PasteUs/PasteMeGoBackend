/*
@File: temporary.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-06-18 17:37  Lucien     1.0         Init
*/
package temporary

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Temporary struct {
	Key       string `gorm:"type:varchar(17);primary_key;index:idx"`
	Lang      string `gorm:"type:varchar(17)"`
	Content   string `gorm:"type:mediumtext"`
	Password  string `gorm:"type:varchar(17)"`
	CreatedAt time.Time
}

// Guarantee there is no such key in temporary table
func Insert(db *gorm.DB, key string, lang string, content string, password string) (string, error) {
	temporary := &Temporary{
		Key:      key,
		Content:  content,
		Lang:     lang,
		Password: password}
	if err := db.Create(&temporary).Error; err != nil {
		return "", err
	}
	return temporary.Key, nil
}

func Query(db *gorm.DB, key string) (Temporary, error) {
	temporary := Temporary{}
	err := db.Find(&temporary, "`key` = ?", key).Error
	return temporary, err
}

func Delete(db *gorm.DB, key string) error {
	if err := db.Where("`key` = ?", key).Delete(Temporary{}).Error; err != nil {
		return err
	}
	return nil
}

func Exist(db *gorm.DB, key string) bool {
	count := uint8(0)
	if err := db.Model(&Temporary{}).Where("`key` = ?", key).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}
