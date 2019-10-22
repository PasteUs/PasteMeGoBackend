/*
@File: model.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-06-23 14:03  Lucien     1.0.0       None
2019-10-20 19:30  Lucien     1.1.0       Add SQLite support
*/
package model

import (
	"fmt"
	"github.com/PasteUs/PasteMeGoBackend/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wonderivan/logger"
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

func formatWithConfig(config config.Config) string {
	database := config.Database
	return format(database.Username, database.Password, "tcp", database.Server, database.Port, database.Database)
}

var db *gorm.DB

func init() {
	var err error
	if config.Data.Database.Type != "mysql" {
		db, err = gorm.Open("sqlite3", "pasteme.db")
		if err != nil {
			logger.Fatal("Connect to SQLite failed: " + err.Error())
		} else {
			logger.Info("SQLite connected")
			if config.Data.Debug {
				logger.Warn("Running in debug mode, database execute will be displayed")
				db = db.Debug()
			}
			if !db.HasTable(&Permanent{}) {
				logger.Warn("Table permanents not found, start creating")
				if err := db.CreateTable(&Permanent{}).Error; err != nil {
					logger.Fatal("Create table permanents failed: " + err.Error())
				}
				db.Exec("INSERT INTO `sqlite_sequence` (`name`, `seq`) VALUES ('permanents', 99)")
			}

			if !db.HasTable(&Temporary{}) {
				logger.Warn("Table temporaries not found, start creating")
				if err := db.CreateTable(&Temporary{}).Error; err != nil {
					logger.Fatal("Create table temporaries failed: " + err.Error())
				}
			}
		}
	} else {
		db, err = gorm.Open("mysql", formatWithConfig(config.Data))
		if err != nil {
			logger.Fatal("Connect to MySQL failed: " + err.Error())
		} else {
			logger.Info("MySQL connected")
			if config.Data.Debug {
				logger.Warn("Running in debug mode, database execute will be displayed")
				db = db.Debug()
			}
			if !db.HasTable(&Permanent{}) {
				logger.Warn("Table permanents not found, start creating")
				if err := db.Set(
					"gorm:table_options",
					"ENGINE=Innodb DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=100",
				).CreateTable(&Permanent{}).Error; err != nil {
					logger.Fatal("Create table permanents failed: " + err.Error())
				}
			}

			if !db.HasTable(&Temporary{}) {
				logger.Warn("Table temporaries not found, start creating")
				if err := db.Set(
					"gorm:table_options",
					"ENGINE=Innodb DEFAULT CHARSET=utf8mb4",
				).CreateTable(&Temporary{}).Error; err != nil {
					logger.Fatal("Create table temporaries failed: " + err.Error())
				}
			}
		}
	}
}
