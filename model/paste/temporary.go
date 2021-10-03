package paste

import (
	"github.com/PasteUs/PasteMeGoBackend/model/dao"
	"gorm.io/gorm"
	"time"
)

const (
	OneMonth = 31 * 24 * 60
	MaxCount = 3
)

func init() {
	dao.CreateTable(&Temporary{})
}

// Temporary 临时
type Temporary struct {
	*AbstractPaste        // 公有字段
	ExpireMinute   uint64 `json:"expire_minute"` // 过期时间
	ExpireCount    uint64 `json:"expire_count"` // 过期的次数
}

// Save 成员函数，保存
func (paste *Temporary) Save() error {
	paste.Key = Generator(8, true, &paste)
	paste.Password = hash(paste.Password)
	return dao.DB.Create(&paste).Error
}

// Delete 成员函数，删除
func (paste *Temporary) Delete() error {
	return dao.DB.Delete(&paste).Error
}

// Get 成员函数，查看
func (paste *Temporary) Get(password string) error {
	err := dao.DB.Transaction(func(tx *gorm.DB) error {
		if e := tx.Take(&paste).Error; e != nil {
			return e
		}

		duration := time.Minute * time.Duration(paste.ExpireMinute)
		if time.Now().After(paste.CreatedAt.Add(duration)) {
			if e := tx.Delete(&paste).Error; e != nil {
				return e
			}
			return gorm.ErrRecordNotFound
		}

		if e := paste.checkPassword(password); e != nil {
			return e
		}

		if paste.ExpireCount <= 1 {
			if e := tx.Delete(&paste).Error; e != nil {
				return e
			}
		} else {
			return tx.Model(&paste).Update("expire_count", paste.ExpireCount-1).Error
		}

		return nil
	})
	return err
}


