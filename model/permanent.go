package model

import (
    "github.com/PasteUs/PasteMeGoBackend/model/dao"
    "time"
)

// Permanent 永久
type Permanent struct {
    Key uint64 `gorm:"primary_key"` // 主键:索引
    *AbstractPaste
    // 存储记录的删除时间
    // 删除具有 DeletedAt 字段的记录，它不会从数据库中删除，但只将字段 DeletedAt 设置为当前时间，并在查询时无法找到记录
    DeletedAt *time.Time
}

// Save 成员函数，创建
func (paste *Permanent) Save() error {
    if err := paste.beforeSave(); err != nil {
        return err
    }
    return dao.Connection().Create(&paste).Error
}

// Delete 成员函数，删除
func (paste *Permanent) Delete() error {
    return dao.Connection().Delete(&paste).Error
}

// Get 成员函数，访问
func (paste *Permanent) Get() error {
    return dao.Connection().First(&paste).Error
}
