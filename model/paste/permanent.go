package paste

import (
	"github.com/PasteUs/PasteMeGoBackend/model/dao"
	"gorm.io/gorm"
)

func init() {
	dao.CreateTable(&Permanent{})
}

// Permanent 永久
type Permanent struct {
	*AbstractPaste
	// 存储记录的删除时间
	// 删除具有 DeletedAt 字段的记录，它不会从数据库中删除，但只将字段 DeletedAt 设置为当前时间，并在查询时无法找到记录
	DeletedAt gorm.DeletedAt
}

// Save 成员函数，创建
func (paste *Permanent) Save() error {
	paste.Key = generator(8, false, &paste)
	paste.Password = hash(paste.Password)
	return dao.DB.Create(&paste).Error
}

// Delete 成员函数，删除
func (paste *Permanent) Delete() error {
	return dao.DB.Delete(&paste).Error
}

func (paste *Permanent) Get(password string) error {
	if err := dao.DB.Take(&paste).Error; err != nil {
		return err
	}
	if err := paste.checkPassword(password); err != nil {
		return err
	}
	return nil
}
