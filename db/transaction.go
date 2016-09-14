package db

import (
	"github.com/jinzhu/gorm"
)

func (db *DB) Begin() *gorm.DB {
	if db.tx == nil {
		db.tx = db.db.Begin()
	}
	return db.tx
}

func (db *DB) Rollback() error {
	err := db.tx.Rollback().Error
	db.tx = nil
	return err
}

func (db *DB) Commit() error {
	err := db.tx.Commit().Error
	db.tx = nil
	return err
}
