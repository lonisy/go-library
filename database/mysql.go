package database

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/schema"
    "sync"
)

type MySQLStruct struct {
    db   *gorm.DB
    once sync.Once
}

var MySQL MySQLStruct

func (s *MySQLStruct) DB() *gorm.DB {
    s.once.Do(func() {
        var err error
        dsn := `AppConfig.GetString("mysql-con")`
        s.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
            NamingStrategy: schema.NamingStrategy{
                TablePrefix:   "prefix_",
                SingularTable: false,
            },
        })
        if err != nil {
            panic(err.Error())
        }
    })
    return s.db
}
