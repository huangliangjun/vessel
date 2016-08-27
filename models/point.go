package models

import (
	"time"
)

const (
	// StarPoint point type
	StarPoint = "Start"
	// CheckPoint point type
	CheckPoint = "Check"
	// EndPoint point type
	EndPoint = "End"

	StartPointMark = "$StartPointMark$"
	EndPointMark   = "$EndPointMark$"
)

// Point data
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

//custom set Point's table name to be pipeline_point
func (Point) TableName() string {
	return "pipeline_point"
}

//query pipeline's points data
func (p *Point) Query() ([]*Point, error) {
	engineDb := Db
	points := make([]*Point, 0, 10)
	p.Status = DataValidStatus
	err := engineDb.Find(&points, p).Error
	return points, err
}
