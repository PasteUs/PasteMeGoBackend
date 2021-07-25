package paste

import (
    "errors"
    "github.com/PasteUs/PasteMeGoBackend/config"
    "github.com/PasteUs/PasteMeGoBackend/util/logging"
    "go.uber.org/zap"
    "time"
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
    Key       string `json:"key" gorm:"type:varchar(16);primary_key;unique_index:idx_temporary"` // 主键:索引
    Namespace string `json:"namespace" gorm:"type:varchar(16);unique_index:idx_temporary"`       // 用户名
    *AbstractPaste
    Expiration time.Time // 过期时间
    Quota      uint8     // 访问计数
}

func (Temporary) TableName() string {
    return "temporary"
}

// Save 成员函数，保存
func (paste *Temporary) Save() error {
    if err := paste.beforeSave(); err != nil {
        return err
    }
    if paste.Expiration.IsZero() && paste.Quota == 0 {
        return errors.New("both expiration and quota disabled")
    }
    return db.Create(&paste).Error
}

// Delete 成员函数，删除
func (paste *Temporary) Delete() error {
    return db.Delete(&paste).Error
}

// Get 成员函数，查看
func (paste *Temporary) Get(password string) error {
    if !paste.Expiration.IsZero() && paste.Expiration.Before(time.Now()) { // 如果存在过期时间且已过期
        if err := paste.Delete(); err != nil {
            return nil
        }
    }

    if err := db.Find(&paste).Error; err != nil {
        return err
    }
    if err := paste.checkPassword(password); err != nil {
        return err
    }

    if paste.Quota == 1 { // 如果剩余浏览次数只有 1 次，那么本次浏览完就会归零，故删除
        if err := paste.Delete(); err != nil {
            return nil
        }
    }

    return nil
}

func Exist(key string) bool {
    count := uint8(0)
    db.Model(&Temporary{}).Where("`key` = ?", key).Count(&count)
    return count > 0
}