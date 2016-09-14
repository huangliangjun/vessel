package models

import (
	"time"

	"github.com/containerops/vessel/db"
)

const (
	// StarPoint point type
	StartPoint = "Start"
	// CheckPoint point type
	CheckPoint = "Check"
	// EndPoint point type
	EndPoint = "End"

	// TemporaryPoint point type
	TemporaryPoint = "Temporary"

	// StartPointMark start mark
	StartPointMark = "$StartPointMark$"
	// EndPointMark end mark
	EndPointMark = "$EndPointMark$"
)

// Point data
type Point struct {
	ID         uint64     `json:"id" gorm:"primary_key"`
	PID        uint64     `json:"pid" gorm:"type:int;not null;index"`
	Type       string     `json:"type" binding:"In(Start,Check,End)" gorm:"type:varchar(20)"`
	Triggers   string     `json:"triggers" gorm:"type:varchar(255)"`   // eg : "a,b,c"
	Conditions string     `json:"conditions" gorm:"type:varchar(255)"` // eg : "a,b,c"
	Status     uint       `json:"status" gorm:"type:tinyint;default:0"`
	CreatedAt  *time.Time `json:"created" `
	UpdatedAt  *time.Time `json:"updated"`
	DeletedAt  *time.Time `json:"deleted"`
}

// PointVersion data
type PointVersion struct {
	ID         uint64     `json:"id" gorm:"primary_key"`
	PvID       uint64     `json:"pvid" gorm:"type:int;not null"`
	PointID    uint64     `json:"PointID" gorm:"type:int;not null;index"`
	State      string     `json:"state" gorm:"column:state;type:varchar(20);not null;"`
	Detail     string     `json:"detail" gorm:"type:text;"`
	Status     uint       `json:"status" gorm:"type:tinyint;default:0"`
	Conditions []string   `json:"-" sql:"-"`
	MetaData   *Point     `json:"-" sql:"-"`
	Kind       string     `json:"-" sql:"-"`
	CreatedAt  *time.Time `json:"created" `
	UpdatedAt  *time.Time `json:"updated"`
	DeletedAt  *time.Time `json:"deleted"`
}

// custom set Point's table name is Point in db
func (p *Point) TableName() string {
	return "point"
}

// custom set PointVersion's table name is point_version in db
func (p *PointVersion) TableName() string {
	return "point_version"
}

//query pipeline's points data
func (p *Point) QueryM() ([]*Point, error) {
	points := make([]*Point, 0, 10)
	err := db.Instance.QueryM(p, &points)
	if err != nil {
		return nil, err
	}
	return points, err
}

//create point's version data
func (pv *PointVersion) Create() error {
	if err := db.Instance.Create(pv); err != nil {
		return err
	}
	return db.Instance.Commit()
}

//update point's version data
func (pv *PointVersion) Update() error {
	if err := db.Instance.Update(pv); err != nil {
		return err
	}
	return db.Instance.Commit()
}

//query point's version data
func (pv *PointVersion) QueryM() ([]*PointVersion, error) {
	pvs := make([]*PointVersion, 0, 10)
	if err := db.Instance.QueryM(pv, &pvs); err != nil {
		return nil, err
	}
	return pvs, nil
}

//delete point's version data
func (pv *PointVersion) SoftDelete() error {
	if err := db.Instance.DeleteS(pv); err != nil {
		return err
	}
	return db.Instance.Commit()
}
