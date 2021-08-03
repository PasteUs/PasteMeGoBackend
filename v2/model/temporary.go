package model

import "github.com/PasteUs/PasteMeGoBackend/model/dao"

// Temporary 临时
type Temporary struct {
	Key string `json:"key" gorm:"type:varchar(16);primary_key"` // 主键:索引
	*AbstractPaste
}

// Save 成员函数，保存
func (paste *Temporary) Save() error {
	if err := paste.beforeSave(); err != nil {
		return err
	}
	return dao.DB().Create(&paste).Error
}

// Delete 成员函数，删除
func (paste *Temporary) Delete() error {
	return dao.DB().Delete(&paste).Error
}

// Get 成员函数，查看
func (paste *Temporary) Get() error {
	return dao.DB().Find(&paste).Error
}

func Exist(key string) bool {
	count := uint8(0)
	dao.DB().Model(&Temporary{}).Where("`key` = ?", key).Count(&count)
	return count > 0
}