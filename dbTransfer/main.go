/*
@File: main.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

Transform database from version 2.x to version 3.x

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-07-01 23:39  Lucien     1.0         Init
*/
package main

import (
	"database/sql"
	"fmt"
	"github.com/LucienShui/PasteMeBackend/dbTransfer/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wonderivan/logger"
	"html"
	"time"
)

var (
	username = "pasteme_cn"
	password = "password"
	network  = "tcp"
	server   = "web"
	port     = 3306
	database = "pasteme_cn"
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

var permCount, tempCount uint64 = 0, 0

func main() {
	if db, err := sql.Open("mysql", format(username, password, network, server, port, database)); err != nil {
		logger.Fatal("Connect to MySQL failed: " + err.Error())
	} else {
		fixAutoIncrement(db)

		start := time.Now()
		permanent(db)
		timePerm := time.Since(start)
		logger.Info("Permanent finished: ", timePerm)

		start = time.Now()
		temporary(db)
		timeTemp := time.Since(start)
		logger.Info("Temporary finished: ", timeTemp)

		logger.Info("=====================================")
		logger.Info(fmt.Sprintf("%d records total, cost: ", permCount+tempCount), timePerm+timeTemp)
		logger.Info(fmt.Sprintf("%d permanents cost: ", permCount), timePerm)
		logger.Info(fmt.Sprintf("%d temporaries cost: ", tempCount), timeTemp)
	}
}

func temporary(db *sql.DB) {
	if rows, err := db.Query("SELECT `key`, `type`, `text`, `passwd` FROM `temp`"); err != nil {
		logger.Fatal("MySQL query failed: " + err.Error())
	} else {
		for rows.Next() {
			object, lang, passwd := model.Temporary{}, sql.NullString{}, sql.NullString{}
			if err := rows.Scan(&object.Key, &lang, &object.Content, &passwd); err != nil {
				logger.Fatal("Scan error: " + err.Error())
			} else {
				if passwd.Valid {
					object.Password = passwd.String
				} else {
					object.Password = ""
				}

				if lang.Valid {
					object.Lang = lang.String
				} else {
					object.Lang = "plain"
				}

				if object.Content == "" {
					object.Content = " "
				}

				object.Content = html.UnescapeString(object.Content)

				if err := object.Save(); err != nil {
					logger.Fatal(fmt.Sprintf("Paste %s save failed: %s", object.Key, err.Error()))
				} else {
					logger.Debug(fmt.Sprintf("Paste %s save successful", object.Key))
					tempCount++
				}
			}
		}
	}
}

func permanent(db *sql.DB) {
	for i := 0; i < 10; i++ {
		if rows, err := db.Query(fmt.Sprintf("SELECT `key`, `type`, `text`, `passwd` FROM `perm%d`", i)); err != nil {
			logger.Fatal("MySQL query failed: " + err.Error())
		} else {
			for rows.Next() {
				object, lang, passwd := model.Permanent{}, sql.NullString{}, sql.NullString{}
				if err := rows.Scan(&object.Key, &lang, &object.Content, &passwd); err != nil {
					logger.Fatal("Scan error: " + err.Error())
				} else {
					if passwd.Valid {
						object.Password = passwd.String
					} else {
						object.Password = ""
					}

					if lang.Valid {
						object.Lang = lang.String
					} else {
						object.Lang = "plain"
					}

					if object.Content == "" {
						object.Content = " "
					}

					object.Content = html.UnescapeString(object.Content)

					if err := object.Save(); err != nil {
						logger.Fatal(fmt.Sprintf("Paste %d save failed: %s", object.Key, err.Error()))
					} else {
						logger.Debug(fmt.Sprintf("Paste %d save successful", object.Key))
						permCount++
					}
				}
			}
		}
	}
}

func fixAutoIncrement(db *sql.DB) {
	var autoIncrement uint64
	if err := db.QueryRow("SELECT `id` FROM `id`").Scan(&autoIncrement); err != nil {
		logger.Fatal("MySQL query from old table: `id` failed: " + err.Error())
	}
	if err := model.FixAutoIncrement(autoIncrement); err != nil {
		logger.Fatal("MySQL alter new table: `permanents`.`AUTO_INCREMENT` failed: " + err.Error())
	}
}
