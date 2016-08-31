package models

import (
	"fmt"
	"time"
)

// PipelineTemplate template for request data
type PipelineTemplate struct {
	Kind       string    `json:"kind" binding:"In(CCloud)"`
	APIVersion string    `json:"apiVersion" binding:"In(v1)"`
	MetaData   *Pipeline `json:"metadata" binding:"Required"`
}

// Pipeline pipeline data
type Pipeline struct {
	ID        uint64     `json:"id" gorm:"primary_key"`
	Namespace string     `json:"namespace" binding:"Required" gorm:"type:varchar(20);not null;unique_index:idxs_namespace_name"`
	Name      string     `json:"name" binding:"Required" gorm:"type:varchar(20);not null;unique_index:idxs_namespace_name"`
	Timeout   uint64     `json:"timeout" gorm:"type:int;"`
	Stages    []*Stage   `json:"stages" binding:"Required" gorm:"-"`
	Points    []*Point   `json:"points" binding:"Required" gorm:"-"`
	Status    uint       `json:"status" gorm:"type:tinyint;default:0"`
	CreatedAt *time.Time `json:"created" `
	UpdatedAt *time.Time `json:"updated"`
	DeletedAt *time.Time `json:"deleted"`
}

// PipelineVersion data
type PipelineVersion struct {
	ID        uint64     `json:"id" gorm:"primary_key"`
	PID       uint64     `json:"Pid" gorm:"type:int;not null;index"`
	State     string     `json:"state" gorm:"column:state;type:varchar(20)"`
	Detail    string     `json:"detail" gorm:"type:text;"`
	MetaData  *Pipeline  `json:"-" sql:"-"`
	Status    uint       `json:"status" gorm:"type:tinyint;default:0"`
	CreatedAt *time.Time `json:"created" `
	UpdatedAt *time.Time `json:"updated"`
	DeletedAt *time.Time `json:"deleted"`
}

// PipelineResult data
type PipelineResult struct {
	PID       uint64 `json:"pid"`
	PvID      uint64 `json:"pvid"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Detail    string `json:"detail"`
}

// custom set pipeline's table name is pipeline_version in db
func (p *Pipeline) TableName() string {
	return "pipeline"
}

//custom set PipelineVersion's table name is pipeline_version in db
func (PipelineVersion) TableName() string {
	return "pipeline_version"
}

//add pipeline data
func (p *Pipeline) Add() error {
	engineDb := db
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
		point.PID = p.ID
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
		stage.PID = p.ID
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
	engineDb := db
	var err error
	//query pipeline data
	err = engineDb.First(p, &Pipeline{Name: p.Name, Namespace: p.Namespace}).Error
	if err != nil {
		return err
	}

	//query pipeline's stages data
	stage := &Stage{PID: p.ID}
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
	point := &Point{PID: p.ID}
	points, err := point.Query()

	if err != nil {
		return err
	}
	p.Points = points

	return nil
}

//delete pipeline data
func (p *Pipeline) Delete() error {
	engineDb := db
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
		PID:    p.ID,
		Status: DataInValidStatus,
	}
	err = tx.Model(&Stage{}).Update(stage).Error
	if err != nil {
		//rollback transaction
		tx.Rollback()
	}

	//modify pipeline'point status
	point := &Point{
		PID:    p.ID,
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
	engineDb := db
	return engineDb.Model(p).Update(p).Error
}

//check pipeline exist
func (p *Pipeline) CheckExist() error {
	engineDb := db
	err := engineDb.First(p, &Pipeline{Name: p.Name, Namespace: p.Namespace}).Error
	fmt.Println(err)
	if err == nil && p.ID > 0 {
		return ErrHasExist
	}
	if err.Error() == ErrNotExist.Error() {
		return nil
	}

	return err
}

//add pipeline version data
func (pv *PipelineVersion) Add() error {
	engineDb := db
	return engineDb.Create(pv).Error
}

//update pipeline version data
func (pv *PipelineVersion) Update() error {
	engineDb := db
	return engineDb.Model(pv).Update(pv).Error
}

//query pipeline version data by pid
func (pv *PipelineVersion) QueryOne() error {
	engineDb := db
	return engineDb.First(pv).Error
}
