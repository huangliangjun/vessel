package db

import (
	//"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type DB struct {
	db *gorm.DB
	tx *gorm.DB
}

var Instance *DB

func init() {
	Instance = &DB{}
}

func InitDB(driver, user, passwd, URI, databaseName string) error {
	if driver == "mysql" {

	} else {
		fmt.Errorf("only suport mysql driver now.")
	}
	//open db
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=Local&parseTime=true", user, passwd, URI, databaseName)
	vesselDb, err := gorm.Open(driver, dsn)
	if err != nil {
		fmt.Errorf("open db error : ", err.Error())
	}
	vesselDb.LogMode(true)
	Instance.db = vesselDb
	Instance.InitConnectedPool()
	return nil
}

func (db *DB) InitConnectedPool() {
	// set connect pool
	db.db.DB().SetMaxIdleConns(10)
	db.db.DB().SetMaxOpenConns(100)
}
