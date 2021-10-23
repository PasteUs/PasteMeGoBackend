package dao

import (
	"fmt"
	"github.com/PasteUs/PasteMeGoBackend/common/config"
	"github.com/PasteUs/PasteMeGoBackend/common/flag"
	"github.com/PasteUs/PasteMeGoBackend/common/logging"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"reflect"
)

func format(
	username string,
	password string,
	network string,
	server string,
	port uint16,
	database string) string {
	return fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, network, server, port, database)
}

func formatWithConfig(database config.Database) string {
	return format(database.Username, database.Password, "tcp", database.Server, database.Port, database.Database)
}

var DB *gorm.DB

func init() {
	var (
		err        error
		gormConfig = &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		}
	)
	if config.Config.Database.Type != "mysql" {
		sqlitePath := config.Config.Database.Database
		pwd, _ := os.Getwd()
		logging.Info("using sqlite", zap.String("database_type", config.Config.Database.Type), zap.String("work_dir", pwd))
		if DB, err = gorm.Open(sqlite.Open(sqlitePath), gormConfig); err != nil {
			logging.Panic("sqlite connect failed", zap.String("sqlite_path", sqlitePath), zap.Error(err))
			return
		}
		logging.Info("sqlite connect success", zap.String("sqlite_path", sqlitePath))
	} else {
		if DB, err = gorm.Open(mysql.Open(formatWithConfig(config.Config.Database)), gormConfig); err != nil {
			logging.Panic("connect to mysql failed", zap.Error(err))
			return
		}
		logging.Info("mysql connected")
	}
	if flag.Debug {
		logging.Warn("running in debug mode, database execute will be displayed")
		DB = DB.Debug()
	}
}

func getTableName(object interface{}) string {
	var typeName string
	if t := reflect.TypeOf(object); t.Kind() == reflect.Ptr {
		typeName = t.Elem().Name()
	} else {
		typeName = t.Name()
	}
	return DB.NamingStrategy.TableName(typeName)
}

func CreateTable(object interface{}) {
	db := DB
	if config.Config.Database.Type == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8mb4")
	}
	migrator := db.Migrator()
	if !migrator.HasTable(object) {
		tableName := zap.String("table_name", getTableName(object))
		logging.Warn("Table not found, start creating", tableName)

		if err := migrator.CreateTable(object); err != nil {
			logging.Panic("Create table failed", tableName, zap.Error(err))
		}
	}
}
