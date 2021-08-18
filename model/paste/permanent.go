package paste

import (
	"github.com/PasteUs/PasteMeGoBackend/model/dao"
	"time"
)

func init() {
	dao.CreateTable(&Permanent{})
	/*
	if !dao.DB.HasTable(&Permanent{}) {
		var err error = nil
		tableName := zap.String("table_name", Permanent{}.TableName())
		logging.Warn("Table not found, start creating", tableName)

		if config.Config.Database.Type != "mysql" {
			err = dao.DB.CreateTable(&Permanent{}).Error
			dao.DB.Exec(fmt.Sprintf("INSERT INTO `sqlite_sequence` (`name`, `seq`) VALUES ('%s', 99)", Permanent{}.TableName()))
		} else {
			err = dao.DB.Set(
				"gorm:table_options",
				"ENGINE=Innodb DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=100",
			).CreateTable(&Permanent{}).Error
		}
		if err != nil {
			logging.Panic("Create table failed", tableName, zap.String("err", err.Error()))
		}
	}
	 */
}

// Permanent 永久
type Permanent struct {
	*AbstractPaste
	// 存储记录的删除时间
	// 删除具有 DeletedAt 字段的记录，它不会从数据库中删除，但只将字段 DeletedAt 设置为当前时间，并在查询时无法找到记录
	DeletedAt *time.Time
}

func (Permanent) TableName() string {
	return "permanent"
}

// Save 成员函数，创建
func (paste *Permanent) Save() error {
	return dao.DB.Create(&paste).Error
}

// Delete 成员函数，删除
func (paste *Permanent) Delete() error {
	return dao.DB.Delete(&paste).Error
}

func (paste *Permanent) Get(password string) error {
	if err := dao.DB.Find(&paste).Error; err != nil {
		return err
	}
	if err := paste.checkPassword(password); err != nil {
		return err
	}
	return nil
}
