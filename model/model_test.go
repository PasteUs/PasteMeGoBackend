/*
@File: model.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-11-20 16:27  bilibili     1.0.0       rewrite
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

func formattest(
	username string,
	password string,
	network string,
	server string,
	port uint16,
	database string) string {
	return fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=True&loc=Local",
		username, password, network, server, port, database)
}

func formatWithConfigtest(config config.Config) string {
	database := config.Database
	return formattest(database.Username, database.Password, "tcp", database.Server, database.Port, database.Database)
}

var testdb *gorm.DB

func Inittest(config config.Config) {
	var err error
	if config.Database.Type != "mysql" {
		testdb, err = gorm.Open("sqlite3", "pasteme.db")
		if err != nil {
			logger.Fatal("Connect to SQLite failed: " + err.Error())
		} else {
			logger.Info("SQLite connected")
			if config.Debug {
				logger.Warn("Running in debug mode, database execute will be displayed")
				testdb = testdb.Debug()
			}
			if !testdb.HasTable(&PermanentUser{}) {
				logger.Warn("Table permanents not found, start creating")
				if err := testdb.CreateTable(&PermanentUser{}).Error; err != nil {
					logger.Fatal("Create table permanents failed: " + err.Error())
				}
				testdb.Exec("INSERT INTO `sqlite_sequence` (`name`, `seq`) VALUES ('permanents', 99)")
			}

			if !testdb.HasTable(&TemporaryUser{}) {
				logger.Warn("Table temporaries not found, start creating")
				if err := testdb.CreateTable(&TemporaryUser{}).Error; err != nil {
					logger.Fatal("Create table temporaries failed: " + err.Error())
				}
			}
		}
	} else {
		testdb, err = gorm.Open("mysql", formatWithConfigtest(config))
		if err != nil {
			logger.Fatal("Connect to MySQL failed: " + err.Error())
		} else {
			logger.Info("MySQL connected")
			if config.Debug {
				logger.Warn("Running in debug mode, database execute will be displayed")
				testdb = testdb.Debug()
			}
			if !testdb.HasTable(&PermanentUser{}) {
				logger.Warn("Table permanents not found, start creating")
				if err := testdb.Set(
					"gorm:table_options",
					"ENGINE=Innodb DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=100",
				).CreateTable(&PermanentUser{}).Error; err != nil {
					logger.Fatal("Create table permanents failed: " + err.Error())
				}
			}
			if !testdb.HasTable(&TemporaryUser{}) {
				logger.Warn("Table temporaries not found, start creating")
				if err := testdb.Set(
					"gorm:table_options",
					"ENGINE=Innodb DEFAULT CHARSET=utf8mb4",
				).CreateTable(&TemporaryUser{}).Error; err != nil {
					logger.Fatal("Create table temporaries failed: " + err.Error())
				}
			}
		}
	}
}

func TestMain(){}