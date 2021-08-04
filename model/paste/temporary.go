package paste

import (
	"github.com/PasteUs/PasteMeGoBackend/config"
	"github.com/PasteUs/PasteMeGoBackend/logging"
	"github.com/PasteUs/PasteMeGoBackend/model/dao"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"time"
)

const (
	EnumTime  = "time"
	EnumCount = "count"
	OneMonth  = 31 * 24 * 60
	MaxCount  = 3
)

func initTemporary() {
	if !dao.DB().HasTable(&Temporary{}) {
		var err error = nil
		tableName := zap.String("table_name", Temporary{}.TableName())
		logging.Warn("Table not found, start creating", tableName)

		if config.Config.Database.Type != "mysql" {
			err = dao.DB().CreateTable(&Temporary{}).Error
		} else {
			err = dao.DB().Set(
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
	Key            string `json:"key" gorm:"type:varchar(16);primary_key"` // 主键:索引
	*AbstractPaste        // 公有字段
	ExpireType     string // 过期类型
	Expiration     uint64 // 过期的数据
}

func (Temporary) TableName() string {
	return "temporary"
}

func (paste *Temporary) GetKey() string {
	return paste.Key
}

func (paste *Temporary) GetNamespace() string {
	return paste.Namespace
}

// Save 成员函数，保存
func (paste *Temporary) Save() error {
	return dao.DB().Create(&paste).Error
}

// Delete 成员函数，删除
func (paste *Temporary) Delete() error {
	return dao.DB().Delete(&paste).Error
}

// Get 成员函数，查看
func (paste *Temporary) Get(password string) error {
	err := dao.DB().Transaction(func(tx *gorm.DB) error {
		if e := tx.Find(&paste).Error; e != nil {
			return e
		}

		if paste.ExpireType == EnumTime {
			duration := time.Minute * time.Duration(paste.Expiration)
			if time.Now().After(paste.CreatedAt.Add(duration)) {
				if e := tx.Delete(&paste).Error; e != nil {
					return e
				}
				return gorm.ErrRecordNotFound
			}
		}

		if e := paste.checkPassword(password); e != nil {
			return e
		}

		if paste.ExpireType == EnumCount {
			if paste.Expiration <= 1 {
				if e := tx.Delete(&paste).Error; e != nil {
					return e
				}
			} else {
				return tx.Model(&paste).Update("expiration", paste.Expiration-1).Error
			}
		}
		return nil
	})
	return err
}

func Exist(key string) bool {
	count := uint8(0)
	dao.DB().Model(&Temporary{}).Where("`key` = ?", key).Count(&count)
	return count > 0
}
