/*
@File: model.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-06-23 14:03  Lucien     1.0         None
*/
package model

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"os"
)

var (
	username = "username"
	password = "password"
	network  = "tcp"
	server   = "mysql"
	port     = 3306
	database = "pasteme"
)

func format(
	username string,
	password string,
	network string,
	server string,
	port int,
	database string) string {
	return fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=True&loc=Local", username, password, network, server, port, database)
}

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", format(username, password, network, server, port, database))
	if err != nil {
		panic(err)
	}
	if os.Getenv("GIN_MODE") != "release" {
		db = db.Debug()
	}
	if !db.HasTable(&Permanent{}) {
		if err := db.Set(
			"gorm:table_options",
			"ENGINE=Innodb DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=100",
		).CreateTable(&Permanent{}).Error; err != nil {
			panic(err)
		}
	}

	if !db.HasTable(&Temporary{}) {
		if err := db.Set(
			"gorm:table_options",
			"ENGINE=Innodb DEFAULT CHARSET=utf8mb4",
		).CreateTable(&Temporary{}).Error; err != nil {
			panic(err)
		}
	}
}
