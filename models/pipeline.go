package models

import (
	"errors"
	//	"fmt"
	"time"

	"github.com/containerops/vessel/db"
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

//create pipeline
func (p *Pipeline) Create() error {
	if err := db.Instance.Create(p); err != nil {
		//rallback transaction
		//db.Instance.Rollback()
		return err
	}
	for _, point := range p.Points {
		//save  pipeline's point data
		point.PID = p.ID
		if err := db.Instance.Create(point); err != nil {
			//db.Instance.Rollback()
			return err
		}
	}

	for _, stage := range p.Stages {
		//save  pipeline's point data
		err := stage.ObjToJson()
		if err != nil {
			//db.Instance.Rollback()
			return err
		}
		stage.PID = p.ID
		if err := db.Instance.Create(stage); err != nil {
			//db.Instance.Rollback()
			return err
		}
	}
	//commit transaction
	return nil
	//return db.Instance.Commit()
}

//Query one pipeline data
func (p *Pipeline) QueryOne() error {
	//query pipeline data
	pipeline := &Pipeline{
		Namespace: p.Namespace,
		Name:      p.Name,
	}
	if _, err := db.Instance.Count(pipeline); err != nil {
		return err
	} else if pipeline.ID <= 0 {
		return errors.New("record not found")
	}
	*p = *pipeline
	//query pipeline's stages data
	stage := &Stage{PID: pipeline.ID}
	stages, err := stage.QueryM()
	if err != nil {
		return err
	}
	p.Stages = stages

	//query pipeline's points data
	point := &Point{PID: pipeline.ID}
	points, err := point.QueryM()
	if err != nil {
		return err
	}
	p.Points = points
	return nil
}

//Query multi pipeline data
func (p *Pipeline) QueryMulti() ([]*Pipeline, error) {
	ps := make([]*Pipeline, 0, 10)
	if err := db.Instance.QueryM(p, &ps); err != nil {
		return nil, err
	}
	return ps, nil
}

//update pipeline data
func (p *Pipeline) Update() error {
	if err := db.Instance.Update(p); err != nil {
		return err
	}
	return nil
}

//delete pipeline data
func (p *Pipeline) SoftDelete() error {
	pipeline := &Pipeline{
		Namespace: p.Namespace,
		Name:      p.Name,
	}
	if _, err := db.Instance.Count(pipeline); err != nil {
		return err
	} else if pipeline.ID <= 0 {
		return errors.New("record not exist")
	}
	//delete pipeline
	if err := db.Instance.DeleteS(pipeline); err != nil {
		return err
	}
	//delete pipeline's stage
	stage := &Stage{
		PID: pipeline.ID,
	}
	if err := db.Instance.DeleteS(stage); err != nil {
		return err
	}

	//delete pipeline's point
	point := &Point{
		PID: pipeline.ID,
	}
	if err := db.Instance.DeleteS(point); err != nil {
		return err
	}

	return nil
}

func (p *Pipeline) CheckIsExist() (bool, error) {
	if _, err := db.Instance.Count(p); err != nil {
		return false, err
	} else if p.ID <= 0 {
		return false, nil
	}
	return true, nil
}

//create pipeline version data
func (pv *PipelineVersion) Create() error {
	if err := db.Instance.Create(pv); err != nil {
		return err
	}
	return nil
}

//update pipeline version data
func (pv *PipelineVersion) Update() error {
	if err := db.Instance.Update(pv); err != nil {
		return err
	}
	return nil
}

//query one pipeline version data
func (pv *PipelineVersion) QueryOne() error {
	if _, err := db.Instance.Count(pv); err != nil {
		return err
	} else if pv.ID <= 0 {
		return errors.New("record not found")
	}
	return nil
}

//query multi pipeline version data
func (pv *PipelineVersion) QueryM() ([]*PipelineVersion, error) {
	pvs := make([]*PipelineVersion, 0, 10)
	if err := db.Instance.QueryM(pv, &pvs); err != nil {
		return nil, err
	}
	return pvs, nil
}

//delete pipeline version data
func (pv *PipelineVersion) SoftDelete() error {
	if err := db.Instance.DeleteS(pv); err != nil {
		return err
	}
	return nil
}
