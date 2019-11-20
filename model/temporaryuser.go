/*
@File: model.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-11-20 16:27  bilibili     1.0.0       rewrite
*/
package model

type TemporaryUser struct {
	Key       string `json:"key" gorm:"type:varchar(16);primary_key"` // 主键:索引
}

func (paste *TemporaryUser) Delete() error {
	return db.Delete(&paste, "`key` = ?", paste.Key).Error
}

// 成员函数，查看
func (paste *TemporaryUser) Get() error {
	return db.Find(&paste, "`key` = ?", paste.Key).Error
}

func TemporaryUserExist(key string) bool {
	count := uint8(0)
	db.Model(&Temporary{}).Where("`key` = ?", key).Count(&count)
	return count > 0
}
