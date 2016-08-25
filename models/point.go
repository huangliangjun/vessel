package models

import (
	"time"
)

const (
	StarPoint  PointKind = "star"
	CheckPoint PointKind = "check"
	EndPoint   PointKind = "end"
)

type PointKind string

type Point struct {
	Id         int64      `json:"id" gorm:"primary_key"`
	Pid        int64      `json:"pid" gorm:"type:bigint;not null;index"`
	Type       string     `json:"type" binding:"In(star,check,end)" gorm:"type:varchar(20)"`
	Triggers   string     `json:"triggers" gorm:"varchar(255)"` // eg : "a,b,c"
	Conditions string     `json:"conditions" gorm:"size:255"`   // eg : "a,b,c"
	Status     uint       `json:"status" gorm:"type:tinyint;default:0"`
	CreatedAt  *time.Time `json:"created" `
	UpdatedAt  *time.Time `json:"updated"`
	DeletedAt  *time.Time `json:"deleted"`
}

//custom table name point
func (Point) TableName() string {
	return "pipeline_point"
}

//query pipeline's points data by pid
func QueryPoints(where map[string]interface{}) (points []*Point, err error) {
	engineDb := Db
	if pid, ok := where["pid"].(int); ok {
		engineDb = engineDb.Where("pid=?", pid)
	}

	err = engineDb.Where("status = ?", DataValidStatus).Find(&points).Error
	return
}
