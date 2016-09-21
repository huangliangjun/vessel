package db

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
)

func (db *DB) RegisterModel(models ...interface{}) error {
	for _, model := range models {
		if !db.db.HasTable(model) {
			if err := db.db.CreateTable(model).Error; err != nil {
				return fmt.Errorf("create table error : ", err.Error())
			}
		}
	}
	return nil
}

func (db *DB) AddUniqueIndex(model interface{}, indexName string, column ...string) error {
	if err := db.db.Model(model).AddUniqueIndex(indexName, column...).Error; err != nil {
		return fmt.Errorf("add unique index error : ", err.Error())
	}
	return nil
}

func (db *DB) AddForeignKey(model interface{}, foreignKeyField, destinationTable, onDelete, onUpdate string) error {
	if err := db.db.Model(model).AddForeignKey(foreignKeyField, destinationTable, onDelete, onUpdate).Error; err != nil {
		return fmt.Errorf("add ForeignKey error : %v", err.Error())
	}
	return nil
}

func (db *DB) Count(value interface{}) (int64, error) {
	var count int64
	if err := db.db.Where(value).Find(value).Count(&count).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

func (db *DB) Create(value interface{}) error {
	tx := db.db.Begin()
	if err := tx.Create(value).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (db *DB) Update(value interface{}) error {
	tx := db.db.Begin()
	if err := tx.Model(value).Update(value).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (db *DB) UpdateField(model interface{}, fieldName string, value interface{}) error {
	return db.db.Model(model).Update(fieldName, value).Error
}

func (db *DB) Save(value interface{}) error {
	tx := db.db.Begin()
	if err := tx.Save(value).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (db *DB) Delete(value interface{}) error {
	tx := db.db.Begin()
	if err := tx.Unscoped().Delete(value).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (db *DB) DeleteS(value interface{}) error {
	tx := db.db.Begin()
	if err := tx.Where(value).Delete(value).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (db *DB) BatchDelete(value interface{}, condition string) error {
	tx := db.db.Begin()
	if err := tx.Unscoped().Where(condition).Delete(value).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (db *DB) BatchDeleteS(value interface{}, condition string) error {
	tx := db.db.Begin()
	if err := tx.Where(condition).Delete(value).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (db *DB) QueryM(condition interface{}, result interface{}) error {
	if err := db.db.Where(condition).Find(result).Error; err != nil {
		return err
	}
	return nil
}

func (db *DB) QueryF(condition interface{}, result interface{}) error {
	var (
		name  string
		value []interface{}
	)
	scope := db.db.NewScope(condition)
	for _, field := range scope.New(condition).Fields() {
		if !field.IsIgnored && !field.IsBlank {
			switch field.Field.Type().Kind() {
			case reflect.String:
				name = name + fmt.Sprintf("%s like ? AND ", field.DBName)
				value = append(value, "%"+field.Field.Interface().(string)+"%")
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Float32, reflect.Float64,
				reflect.Bool:
				name = name + fmt.Sprintf("%s = ? AND ", field.DBName)
				value = append(value, field.Field.Interface())
			default:
				panic("unsuport query type")
			}
		}
	}
	name = strings.TrimRight(name, "AND ")
	if err := db.db.Where(name, value...).Find(result).Error; err != nil {
		return err
	}
	return nil
}

func (db *DB) Raw(models interface{}, sql string, values ...interface{}) error {
	return db.db.Raw(sql, values...).Scan(models).Error
}

func (db *DB) Exec(sql string, values ...interface{}) error {
	tx := db.db.Begin()
	if err := tx.Exec(sql, values...).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
