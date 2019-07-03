/*
@File: model.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-06-23 14:03  Lucien     1.0         None
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/LucienShui/PasteMeBackend/util"
	"github.com/LucienShui/PasteMeBackend/util/convert"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/wonderivan/logger"
	"html"
	"os"
	"strings"
	"time"
)

var (
	username = util.GetEnvOrFatal("PASTEMED_DB_USERNAME")
	password = util.GetEnvOrFatal("PASTEMED_DB_PASSWORD")
	network  = "tcp"
	server   = util.GetEnvOrFatal("PASTEMED_DB_SERVER")
	port     = convert.String2uint(util.GetEnvOrFatal("PASTEMED_DB_PORT"))
	database = util.GetEnvOrFatal("PASTEMED_DB_DATABASE")
)

func format(
	username string,
	password string,
	network string,
	server string,
	port uint64,
	database string) string {
	return fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=True&loc=Local", username, password, network, server, port, database)
}

type Temporary struct {
	Key       string `json:"key" gorm:"type:varchar(16);primary_key;index:idx"`
	Lang      string `json:"lang" gorm:"type:varchar(16)"`
	Content   string `json:"content" gorm:"type:mediumtext"`
	Password  string `json:"password" gorm:"type:varchar(32)"`
	ClientIP  string `gorm:"type:varchar(64)"`
	CreatedAt time.Time
}

func (paste *Temporary) Save(db *gorm.DB) error {
	if paste.Content == "" {
		return errors.New("empty content")
	}
	if paste.Lang == "" {
		return errors.New("empty lang")
	}
	if strings.Contains(paste.Content, "#include") && paste.Lang == "plain" {
		paste.Lang = "cpp"
	}
	return db.Create(&paste).Error
}

type Permanent struct {
	Key       uint64 `gorm:"primary_key;index:idx"`
	Lang      string `json:"lang" gorm:"type:varchar(16)"`
	Content   string `json:"content" gorm:"type:mediumtext"`
	Password  string `json:"password" gorm:"type:varchar(32)"`
	ClientIP  string `gorm:"type:varchar(64)"`
	CreatedAt time.Time
	DeletedAt *time.Time
}

func (paste *Permanent) Save(db *gorm.DB) error {
	if paste.Content == "" {
		return errors.New("empty content")
	}
	if paste.Lang == "" {
		return errors.New("empty lang")
	}
	if strings.Contains(paste.Content, "#include") && paste.Lang == "plain" {
		paste.Lang = "cpp"
	}
	return db.Create(&paste).Error
}

var gormDB *gorm.DB
var sqlDB *sql.DB

func init() {
	var err error
	if sqlDB, err = sql.Open("mysql", format(username, password, network, server, port, database)); err != nil {
		logger.Fatal("Connect to MySQL failed: " + err.Error())
	}

	gormDB, err = gorm.Open("mysql", format(username, password, network, server, port, database))
	if err != nil {
		logger.Fatal("Connect to MySQL failed: " + err.Error())
	} else {
		logger.Info("MySQL connected")
		if os.Getenv("PASTEMED_RUNTIME") == "debug" {
			logger.Warn("Running in debug mode, database execute will be displayed")
			gormDB = gormDB.Debug()
		}
		if !gormDB.HasTable(&Permanent{}) {
			logger.Warn("Table permanents not found, start creating")
			if err := gormDB.Set(
				"gorm:table_options",
				"ENGINE=Innodb DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=100",
			).CreateTable(&Permanent{}).Error; err != nil {
				logger.Fatal("Create table permanents failed: " + err.Error())
			}
		}

		if !gormDB.HasTable(&Temporary{}) {
			logger.Warn("Table temporaries not found, start creating")
			if err := gormDB.Set(
				"gorm:table_options",
				"ENGINE=Innodb DEFAULT CHARSET=utf8mb4",
			).CreateTable(&Temporary{}).Error; err != nil {
				logger.Fatal("Create table temporaries failed: " + err.Error())
			}
		}
	}
}

func fixAutoIncrement(autoIncrement uint64) error {
	return gormDB.Exec(fmt.Sprintf("ALTER TABLE `permanents` AUTO_INCREMENT=%d", autoIncrement)).Error
}

func TransTemporary() uint64 {
	var count uint64 = 0
	if rows, err := sqlDB.Query("SELECT `key`, `type`, `text`, `passwd` FROM `temp`"); err != nil {
		logger.Fatal("MySQL query failed: " + err.Error())
	} else {
		event := gormDB.Begin()
		for rows.Next() {
			object, lang, passwd := Temporary{}, sql.NullString{}, sql.NullString{}
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

				if err := object.Save(event); err != nil {
					logger.Fatal(fmt.Sprintf("Paste %s save failed: %s", object.Key, err.Error()))
				} else {
					logger.Debug(fmt.Sprintf("Paste %s save successful", object.Key))
					count++
				}
			}
		}
		event.Commit()
	}
	return count
}

func TransPermanent() uint64 {
	event := gormDB.Begin()
	var count uint64 = 0
	for i := 0; i < 10; i++ {
		if rows, err := sqlDB.Query(fmt.Sprintf("SELECT `key`, `type`, `text`, `passwd` FROM `perm%d`", i)); err != nil {
			logger.Fatal("MySQL query failed: " + err.Error())
		} else {
			for rows.Next() {
				object, lang, passwd := Permanent{}, sql.NullString{}, sql.NullString{}
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

					if err := object.Save(event); err != nil {
						logger.Fatal(fmt.Sprintf("Paste %d save failed: %s", object.Key, err.Error()))
					} else {
						logger.Debug(fmt.Sprintf("Paste %d save successful", object.Key))
						count++
					}
				}
			}
		}
	}
	event.Commit()
	return count
}

func FixAutoIncrement() {
	var autoIncrement uint64
	if err := sqlDB.QueryRow("SELECT `id` FROM `id`").Scan(&autoIncrement); err != nil {
		logger.Fatal("MySQL query from old table: `id` failed: " + err.Error())
	}
	if err := fixAutoIncrement(autoIncrement); err != nil {
		logger.Fatal("MySQL alter new table: `permanents`.`AUTO_INCREMENT` failed: " + err.Error())
	}
}
