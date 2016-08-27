package models

import (
	//	"encoding/json"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const (
	CCloud PipelineKind = "CCloud"
	V1     APIVersion   = "v1"
)

type PipelineKind string
type APIVersion string

// PipelineSpecTemplate template for request data
type PipelineSpecTemplate struct {
	Kind       PipelineKind `json:"kind" binding:"In(CCloud)"`
	APIVersion APIVersion   `json:"apiVersion" binding:"In(v1)"`
	MetaData   *Pipeline    `json:"metadata" binding:"Required"`
}

// Pipeline pipeline data
type Pipeline struct {
	Id        int64      `json:"id" gorm:"primary_key"`
	Namespace string     `json:"namespace" binding:"Required" gorm:"type:varchar(50);not null;unique_index:idx_namespace_name"`
	Name      string     `json:"name" binding:"Required" gorm:"type:varchar(30);not null;unique_index:idx_namespace_name"`
	Timeout   uint64     `json:"timeout" gorm:"type:int;"`
	Status    uint       `json:"status" gorm:"type:tinyint;default:0"`
	CreatedAt *time.Time `json:"created" `
	UpdatedAt *time.Time `json:"updated"`
	DeletedAt *time.Time `json:"deleted"`
	Stages    []*Stage   `json:"stages" gorm:"-"`
	Points    []*Point   `json:"points" gorm:"-"`
}

//func (Pipeline) TableName() string {
//	return "pipeline"
//}

// PipelineVersion data
type PipelineVersion struct {
	Id            int64      `json:"id" gorm:"primary_key"`
	Pid           int64      `json:"Pid" gorm:"type:bigint;not null;index"`
	Detail        string     `json:"detail" gorm:"type:text;"`
	VersionStatus string     `json:"versionStatus" gorm:"column:versionStatus;type:varchar(20);not null;"`
	Status        uint       `json:"status" gorm:"type:tinyint;default:0"`
	CreatedAt     *time.Time `json:"created" `
	UpdatedAt     *time.Time `json:"updated"`
	DeletedAt     *time.Time `json:"deleted"`
}

// PipelineResult pipeline result
type PipelineResult struct {
	PID       uint   `json:"pid"`
	PvID      uint   `json:"pvid"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

func init() {
	fmt.Println("the Db is ", Db)
	var err error
	if Db == nil {
		dbArgs := "root@tcp(127.0.0.1:3306)/vesseldb?loc=Local&parseTime=True&charset=utf8"
		Db, err = gorm.Open("mysql", dbArgs)

		if err != nil {
			panic(err)
		}
		Db.LogMode(true)
		Db.DB().SetMaxIdleConns(10)
		Db.DB().SetMaxOpenConns(100)
		Db.SingularTable(true)
		if err = Sync(); err != nil {
			panic(err)
		}
	}
}

//custom  set Pipeline's table name to be pipeline
func (Pipeline) TableName() string {
	return "pipeline"
}

//custom set PipelineVersion's table name to be pipeline_version
func (PipelineVersion) TableName() string {
	return "pipeline_version"
}

//add pipeline data
func (p *Pipeline) Add() error {
	engineDb := Db
	//begin transaction
	tx := engineDb.Begin()
	var err error
	//save pipeline data
	if err = tx.Create(p).Error; err != nil {
		//rallback transaction
		tx.Rollback()
		return err
	}

	for _, point := range p.Points {
		//save  pipeline's point data
		point.Pid = p.Id
		if err = tx.Create(point).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, stage := range p.Stages {
		//save  pipeline's point data
		err := stage.ObjToJson()
		if err != nil {
			tx.Rollback()
			return err
		}
		stage.Pid = p.Id
		if err = tx.Create(stage).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	//commit transaction
	tx.Commit()
	return err
}

//Query pipeline data
func (p *Pipeline) QueryOne() error {
	engineDb := Db
	var err error
	//query pipeline data
	err = engineDb.First(p, p).Error
	if err != nil {
		return err
	}

	//query pipeline's stages data
	stage := &Stage{Pid: p.Id}
	stages, err := stage.Query()

	if err != nil {
		return err
	}
	for i, s := range stages {
		err := s.JsonToObj()
		stages[i] = s
		if err != nil {
			return err
		}
	}
	p.Stages = stages

	//query pipeline's points data
	point := &Point{Pid: p.Id}
	points, err := point.Query()

	if err != nil {
		return err
	}
	p.Points = points

	return nil
}

//delete pipeline data
func (p *Pipeline) Delete() error {
	engineDb := Db
	//begin transaction
	tx := engineDb.Begin()

	//modify pipeline status
	p.Status = DataInValidStatus
	err := tx.Model(p).Update(p).Error
	if err != nil {
		//rollback transaction
		tx.Rollback()
	}

	//modify pipeline'stage status
	stage := &Stage{
		Pid:    p.Id,
		Status: DataInValidStatus,
	}
	err = tx.Model(&Stage{}).Update(stage).Error
	if err != nil {
		//rollback transaction
		tx.Rollback()
	}

	//modify pipeline'point status
	point := &Point{
		Pid:    p.Id,
		Status: DataInValidStatus,
	}
	err = tx.Model(&Point{}).Update(point).Error
	if err != nil {
		//rollback transaction
		tx.Rollback()
	}

	//commit transaction
	return tx.Commit().Error
}

//update pipeline data
func (p *Pipeline) Update() error {
	engineDb := Db
	return engineDb.Model(p).Update(p).Error
}

//add pipeline version data
func (pv *PipelineVersion) Add() error {
	engineDb := Db
	return engineDb.Create(pv).Error
}

//update pipeline version data
func (pv *PipelineVersion) Update() error {
	engineDb := Db
	return engineDb.Model(pv).Update(pv).Error
}

//query pipeline version data by pid
func (pv *PipelineVersion) QueryOne() error {
	engineDb := Db
	return engineDb.First(pv).Error
}
