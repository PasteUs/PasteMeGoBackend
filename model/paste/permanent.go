package paste

import (
    "fmt"
    "github.com/PasteUs/PasteMeGoBackend/config"
    "github.com/PasteUs/PasteMeGoBackend/util/convert"
    "github.com/PasteUs/PasteMeGoBackend/util/logging"
    "go.uber.org/zap"
    "time"
)

func init() {
    if !db.HasTable(&Permanent{}) {
        var err error = nil
        tableName := zap.String("table_name", Permanent{}.TableName())
        logging.Warn("Table not found, start creating", tableName)

        if config.Get().Database.Type != "mysql" {
            err = db.CreateTable(&Permanent{}).Error
            db.Exec(fmt.Sprintf("INSERT INTO `sqlite_sequence` (`name`, `seq`) VALUES ('%s', 99)", Permanent{}.TableName()))
        } else {
            err = db.Set(
                "gorm:table_options",
                "ENGINE=Innodb DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=100",
            ).CreateTable(&Permanent{}).Error
        }
        if err != nil {
            logging.Panic("Create table failed", tableName, zap.String("err", err.Error()))
        }
    }
}

// Permanent 永久
type Permanent struct {
    Key       uint64 `gorm:"primary_key;unique_index:idx_permanent"`                       // 主键:索引
    Namespace string `json:"namespace" gorm:"type:varchar(16);unique_index:idx_permanent"` // 用户名
    *AbstractPaste
    // 存储记录的删除时间
    // 删除具有 DeletedAt 字段的记录，它不会从数据库中删除，但只将字段 DeletedAt 设置为当前时间，并在查询时无法找到记录
    DeletedAt *time.Time
}

func (Permanent) TableName() string {
    return "permanent"
}

func (paste *Permanent) GetKey() string {
    return convert.Uint2string(paste.Key)
}

func (paste *Permanent) GetNamespace() string {
    return paste.Namespace
}

// Save 成员函数，创建
func (paste *Permanent) Save() error {
    if err := paste.beforeSave(); err != nil {
        return err
    }
    return db.Create(&paste).Error
}

// Delete 成员函数，删除
func (paste *Permanent) Delete() error {
    return db.Delete(&paste).Error
}

func (paste *Permanent) Get(password string) error {
    if err := db.Find(&paste).Error; err != nil {
        return err
    }
    if err := paste.checkPassword(password); err != nil {
        return err
    }
    return nil
}
