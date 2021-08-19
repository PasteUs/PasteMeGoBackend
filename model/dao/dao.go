package dao

import (
	"fmt"
	"github.com/PasteUs/PasteMeGoBackend/config"
	"github.com/PasteUs/PasteMeGoBackend/flag"
	"github.com/PasteUs/PasteMeGoBackend/logging"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"os"
)

func format(
	username string,
	password string,
	network string,
	server string,
	port uint16,
	database string) string {
	return fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=True&loc=Local",
		username, password, network, server, port, database)
}

func formatWithConfig(database config.Database) string {
	return format(database.Username, database.Password, "tcp", database.Server, database.Port, database.Database)
}

var DB *gorm.DB

func init() {
	var err error
	if config.Config.Database.Type != "mysql" {
		sqlitePath := flag.DataDir + "pasteme.db"
		pwd, _ := os.Getwd()
		logging.Info("using sqlite", zap.String("database_type", config.Config.Database.Type), zap.String("work_dir", pwd))
		if DB, err = gorm.Open("sqlite3", sqlitePath); err != nil {
			logging.Panic("sqlite connect failed", zap.String("sqlite_path", sqlitePath), zap.String("err", err.Error()))
			return
		}
		logging.Info("sqlite connect success", zap.String("sqlite_path", sqlitePath))
	} else {
		if DB, err = gorm.Open("mysql", formatWithConfig(config.Config.Database)); err != nil {
			logging.Panic("connect to mysql failed", zap.String("err", err.Error()))
			return
		}
		logging.Info("mysql connected")
	}
	if flag.Debug {
		logging.Warn("running in debug mode, database execute will be displayed")
		DB = DB.Debug()
	}
}

func CreateTable(object interface{}) {
	if !DB.HasTable(object) {
		var err error = nil
		tableName := zap.String("table_name", DB.NewScope(object).TableName())
		logging.Warn("Table not found, start creating", tableName)

		if config.Config.Database.Type != "mysql" {
			err = DB.CreateTable(object).Error
		} else {
			err = DB.Set(
				"gorm:table_options",
				"ENGINE=Innodb DEFAULT CHARSET=utf8mb4",
			).CreateTable(object).Error
		}

		if err != nil {
			logging.Panic("Create table failed", tableName, zap.String("err", err.Error()))
		}
	}
}
