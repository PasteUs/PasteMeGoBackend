package paste

import (
	"github.com/PasteUs/PasteMeGoBackend/common/logging"
	"github.com/PasteUs/PasteMeGoBackend/model/dao"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

const (
	OneMonth = 31 * 24 * 60 * 60
	MaxCount = 3
)

var nilTime = time.Time{}

func init() {
	dao.CreateTable(&Temporary{})
}

// Temporary 临时
type Temporary struct {
	*AbstractPaste        // 公有字段
	ExpireSecond   uint64 `json:"expire_second"` // 过期时间
	ExpireCount    uint64 `json:"expire_count"` // 过期的次数
}

// Save 成员函数，保存
func (paste *Temporary) Save() error {
	paste.Key = generator(8, true, &paste)
	paste.Password = hash(paste.Password)
	err := dao.DB.Create(&paste).Error
	if err == nil {
		key := paste.Key
		time.AfterFunc(time.Second*time.Duration(paste.ExpireSecond), func() {
			if e := dao.DB.Delete(&Temporary{AbstractPaste: &AbstractPaste{Key: key}}).Error; e != nil {
				logging.Error("delete paste failed", zap.String("key", paste.Key), zap.String("err", e.Error()))
			}
		})
	}
	return err
}

// Delete 成员函数，删除
func (paste *Temporary) Delete() error {
	return dao.DB.Delete(&paste).Error
}

func (paste *Temporary) Expired() bool {
	duration := time.Second * time.Duration(paste.ExpireSecond)
	if time.Now().After(paste.CreatedAt.Add(duration)) {
		return true
	}
	if paste.ExpireCount < 1 {
		return true
	}
	return false
}

// Get 成员函数，查看
func (paste *Temporary) Get(password string) error {
	if err := dao.DB.Transaction(func(tx *gorm.DB) error {
		if e := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&paste).Error; e != nil {
			return e
		}

		if paste.Expired() {
			paste.CreatedAt = nilTime // 通过此字段标记为非法，transaction 生效后再 return error
			return tx.Delete(&paste).Error
		}

		if e := paste.checkPassword(password); e != nil {
			return e
		}

		paste.ExpireCount -= 1

		if paste.Expired() {
			return tx.Delete(&paste).Error
		} else {
			return tx.Save(&paste).Error
		}
	}); err != nil {
		return err
	} else if paste.CreatedAt.IsZero() {
		return gorm.ErrRecordNotFound
	}
	return nil
}
