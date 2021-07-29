package dao

import (
    "fmt"
    "github.com/PasteUs/PasteMeGoBackend/config"
    "github.com/PasteUs/PasteMeGoBackend/flag"
    "github.com/PasteUs/PasteMeGoBackend/util"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
    "go.uber.org/zap"
    "sync"
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

var (
    db   *gorm.DB
    once sync.Once
)

func Init() {
    var err error
    if config.Get().Database.Type != "mysql" {
        sqlitePath := flag.GetArgv().DataDir + "pasteme.db"
        if db, err = gorm.Open("sqlite3", sqlitePath); err != nil {
            util.Panic("sqlite connect failed", zap.String("sqlite_path", sqlitePath), zap.String("err", err.Error()))
            return
        }
        util.Info("sqlite connect success", zap.String("sqlite_path", sqlitePath))
    } else {
        if db, err = gorm.Open("mysql", formatWithConfig(config.Get())); err != nil {
            util.Panic("connect to mysql failed", zap.String("err", err.Error()))
            return
        }
        util.Info("mysql connected")
    }
    if flag.GetArgv().Debug {
        util.Warn("running in debug mode, database execute will be displayed")
        db = db.Debug()
    }
}

func DB() *gorm.DB {
    once.Do(Init)
    return db
}
