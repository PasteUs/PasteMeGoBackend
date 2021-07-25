package paste

import (
    "github.com/PasteUs/PasteMeGoBackend/config"
    "github.com/PasteUs/PasteMeGoBackend/util/logging"
    "go.uber.org/zap"
)

func init() {
    if !db.HasTable(&Temporary{}) {
        var err error = nil
        tableName := zap.String("table_name", Temporary{}.TableName())
        logging.Warn("Table not found, start creating", tableName)

        if config.Get().Database.Type != "mysql" {
            err = db.CreateTable(&Temporary{}).Error
        } else {
            err = db.Set(
                "gorm:table_options",
                "ENGINE=Innodb DEFAULT CHARSET=utf8mb4",
            ).CreateTable(&Temporary{}).Error
        }

        if err != nil {
            logging.Panic("Create table failed", tableName, zap.String("err", err.Error()))
        }
    }
}

// Temporary 临时
type Temporary struct {
    Key string `json:"key" gorm:"type:varchar(16);primary_key"` // 主键:索引
    *AbstractPaste
}

func (Temporary) TableName() string {
    return "temporary"
}

// Save 成员函数，保存
func (paste *Temporary) Save() error {
    if err := paste.beforeSave(); err != nil {
        return err
    }
    return db.Create(&paste).Error
}

// Delete 成员函数，删除
func (paste *Temporary) Delete() error {
    return db.Delete(&paste).Error
}

// Get 成员函数，查看
func (paste *Temporary) Get() error {
    return db.Find(&paste).Error
}

func Exist(key string) bool {
    count := uint8(0)
    db.Model(&Temporary{}).Where("`key` = ?", key).Count(&count)
    return count > 0
}
