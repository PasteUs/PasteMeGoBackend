package user

import (
	"github.com/PasteUs/PasteMeGoBackend/common/config"
	"github.com/PasteUs/PasteMeGoBackend/common/logging"
	"github.com/PasteUs/PasteMeGoBackend/model/dao"
	"go.uber.org/zap"
)

func init() {
	if !dao.DB.HasTable(&User{}) {
		var err error = nil
		tableName := zap.String("table_name", User{}.TableName())
		logging.Warn("Table not found, start creating", tableName)

		if config.Config.Database.Type != "mysql" {
			err = dao.DB.CreateTable(&User{}).Error
		} else {
			err = dao.DB.Set(
				"gorm:table_options",
				"ENGINE=Innodb DEFAULT CHARSET=utf8mb4",
			).CreateTable(&User{}).Error
		}

		if err != nil {
			logging.Panic("Create table failed", tableName, zap.String("err", err.Error()))
		}
	}
}

type User struct {
	Username string `json:"username" json:"username" binding:"required"`
	Password string `json:"password" json:"password" binding:"required"`
}

func (User) TableName() string {
	return "user"
}
