package dao

import (
    "fmt"
    "github.com/PasteUs/PasteMeGoBackend/config"
    "github.com/PasteUs/PasteMeGoBackend/flag"
    "github.com/PasteUs/PasteMeGoBackend/util/logging"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
    "go.uber.org/zap"
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
    if config.Get().Database.Type != "mysql" {
        sqlitePath := flag.DataDir + "pasteme.db"
        db, err = gorm.Open("sqlite3", sqlitePath)
        if err != nil {
            logging.Panic("SQLite connect failed", zap.String("sqlite_path", sqlitePath), zap.String("err", err.Error()))
        } else {
            logging.Info("SQLite connect success", zap.String("sqlite_path", sqlitePath))
            if flag.Debug {
                logging.Warn("Running in debug mode, database execute will be displayed")
                db = db.Debug()
            }
        }
    } else {
        db, err = gorm.Open("mysql", formatWithConfig(config.Get()))
        if err != nil {
            logging.Panic("Connect to MySQL failed", zap.String("err", err.Error()))
        } else {
            logging.Info("MySQL connected")
            if flag.Debug {
                logging.Warn("Running in debug mode, database execute will be displayed")
                db = db.Debug()
            }
        }
    }
}

func Connection() *gorm.DB {
    return db
}
