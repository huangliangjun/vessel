package models

import (
	//	"errors"
	"fmt"
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
	PID        uint64     `json:"pID" gorm:"not null;index"`
	Type       string     `json:"type" binding:"In(Start,Check,End)" gorm:"type:varchar(32)"`
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
	PvID       uint64     `json:"pvID" gorm:"not null;"`
	PointID    uint64     `json:"pointID" gorm:"not null;"`
	State      string     `json:"state" gorm:"type:varchar(32);not null;"`
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
func (pv *PointVersion) TableName() string {
	return "point_version"
}

func (p *Point) AddForeignKey() error {
	if err := db.Instance.AddForeignKey(p, "p_id", "pipeline(id)", "CASCADE", "NO ACTION"); err != nil {
		return fmt.Errorf("create foreign key p_id error: %v", err.Error())
	}
	return nil
}

func (pv *PointVersion) AddForeignKey() error {
	if err := db.Instance.AddForeignKey(pv, "pv_id", "pipeline_version(id)", "CASCADE", "NO ACTION"); err != nil {
		return fmt.Errorf("create foreign key pv_id error: %v", err.Error())
	}
	return nil
}

func (p *PointVersion) AddUniqueIndex() error {
	if err := db.Instance.AddUniqueIndex(p, "idxs_pvid_pointid", "pv_id", "point_id"); err != nil {
		return fmt.Errorf("create unique index idxs_pvid_pointid error: %v", err.Error())
	}
	return nil
}

//check point record is exist
func (p *Point) IsExist() (bool, error) {
	if _, err := db.Instance.Count(p); err != nil {
		return false, err
	} else if p.ID <= 0 {
		return false, nil
	}
	return true, nil
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
	return nil
}

//update point's version data
func (pv *PointVersion) Update() error {
	pointVersion := &PointVersion{
		PvID:    pv.PvID,
		PointID: pv.PointID,
	}
	is, err := pointVersion.IsExist()
	if err != nil {
		return err
	} else if err == nil && is == false {
		return fmt.Errorf("record not exist")
	}
	//	if _, err := db.Instance.Count(pointVersion); err != nil {
	//		return err
	//	} else if pointVersion.ID <= 0 {
	//		return errors.New("record not exist")
	//	}
	pv.ID = pointVersion.ID
	if err := db.Instance.Update(pv); err != nil {
		return err
	}
	return nil
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
	return nil
}

//check pointVersion record is exist
func (p *PointVersion) IsExist() (bool, error) {
	if _, err := db.Instance.Count(p); err != nil {
		return false, err
	} else if p.ID <= 0 {
		return false, nil
	}
	return true, nil
}
