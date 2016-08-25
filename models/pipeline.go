package models

import (
	"encoding/json"
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

//custom table name pipeline
func (Pipeline) TableName() string {
	return "pipeline"
}

//custom table name pipeline_version
func (PipelineVersion) TableName() string {
	return "pipeline_version"
}

//add pipeline data
func AddPipeline(pipelineSpecTemplate *PipelineSpecTemplate) (id int64, err error) {
	engineDb := Db
	//begin transaction
	tx := engineDb.Begin()

	pipeline := Pipeline{
		Namespace: pipelineSpecTemplate.MetaData.Namespace,
		Name:      pipelineSpecTemplate.MetaData.Name,
		Timeout:   pipelineSpecTemplate.MetaData.Timeout,
	}

	//save pipeline data
	if err = tx.Create(&pipeline).Error; err != nil {
		//rallback transaction
		tx.Rollback()
		return
	}
	if pipelineSpecTemplate.MetaData.Points != nil {
		for _, value := range pipelineSpecTemplate.MetaData.Points {
			point := Point{
				Pid:        pipeline.Id,
				Type:       value.Type,
				Triggers:   value.Triggers,
				Conditions: value.Conditions,
			}
			//save  pipeline's point data
			if err = tx.Create(&point).Error; err != nil {
				tx.Rollback()
				return
			}
		}
	}

	if pipelineSpecTemplate.MetaData.Stages != nil {
		for _, value := range pipelineSpecTemplate.MetaData.Stages {
			bsArtifacts, err := json.Marshal(value.Artifacts)
			if err != nil {
				tx.Rollback()
				return -1, err
			}
			artifactsJson := string(bsArtifacts)
			bsVolumes, err := json.Marshal(value.Volumes)
			if err != nil {
				tx.Rollback()
				return -1, err
			}
			volumesJson := string(bsVolumes)
			stage := Stage{
				Pid:           pipeline.Id,
				Name:          value.Name,
				Type:          value.Type,
				Dependencies:  value.Dependencies,
				ArtifactsJson: artifactsJson,
				VolumesJson:   volumesJson,
			}
			//save pipeline's stage data
			if err = tx.Create(&stage).Error; err != nil {
				tx.Rollback()
				return -1, err
			}
		}
	}
	//commit transaction
	tx.Commit()
	return pipeline.Id, err
}

//Query pipeline data
func QueryPipeline(where map[string]interface{}) (*PipelineSpecTemplate, error) {
	engineDb := Db
	var pipeline Pipeline
	if id, ok := where["id"].(int); ok {
		engineDb = engineDb.Where("id=?", id)
	}
	if namespace, ok := where["namespace"].(string); ok {
		engineDb = engineDb.Where("namespace=?", namespace)
	}
	if name, ok := where["name"].(string); ok {
		engineDb = engineDb.Where("name=?", name)
	}
	//query pipeline data
	err := engineDb.Where("status = ?", DataValidStatus).First(&pipeline).Error
	if err != nil {
		return nil, err
	}
	stageWhere := map[string]interface{}{
		"pid": pipeline.Id,
	}
	//query pipeline's stages data
	stages, err := QueryStages(stageWhere)
	if err != nil {
		return nil, err
	}
	for i, value := range stages {
		var artifacts []*Artifact
		var volumes []*Volume
		json.Unmarshal([]byte(value.ArtifactsJson), &artifacts)
		stages[i].Artifacts = artifacts
		json.Unmarshal([]byte(value.VolumesJson), &volumes)
		stages[i].Volumes = volumes
	}
	pipeline.Stages = stages

	pointWhere := map[string]interface{}{
		"pid": pipeline.Id,
	}
	//query pipeline's points data
	points, err := QueryPoints(pointWhere)
	if err != nil {
		return nil, err
	}
	pipeline.Points = points
	pipelineSpecTemplate := new(PipelineSpecTemplate)
	pipelineSpecTemplate.MetaData = &pipeline

	return pipelineSpecTemplate, nil
}

//delete pipeline data
func DeletePipeline(id int64) (err error) {
	engineDb := Db
	//begin transaction
	tx := engineDb.Begin()
	//modify pipeline status
	err = tx.Table("pipeline").Where("id=?", id).Update(&Pipeline{Status: DataInValidStatus}).Error
	if err != nil {
		//rollback transaction
		tx.Rollback()
	}
	//modify pipeline'stage status
	err = tx.Table("pipeline_stage").Where("pid=?", id).Update(&Stage{Status: DataInValidStatus}).Error
	if err != nil {
		tx.Rollback()
	}
	//modify pipeline'point status
	err = tx.Table("pipeline_point").Where("pid=?", id).Update(&Point{Status: DataInValidStatus}).Error
	if err != nil {
		tx.Rollback()
	}
	//commit transaction
	err = tx.Commit().Error

	return
}

//update pipeline data
func UpdatePipeline(pipelineSpecTemplate *PipelineSpecTemplate) (err error) {
	return
}

//save pipeline version data
func AddPipelineVersion(pipelineVersion *PipelineVersion) (id int64, err error) {
	engineDb := Db
	err = engineDb.Create(pipelineVersion).Error
	return pipelineVersion.Id, err
}

//update pipeline version data
func UpdatePipelineVersion(pipelineVersion *PipelineVersion) (err error) {
	engineDb := Db
	err = engineDb.Table("pipeline_version").Where("pid = ?", pipelineVersion.Pid).Update(
		&PipelineVersion{
			VersionStatus: pipelineVersion.VersionStatus}).Error
	return
}

//query pipeline version data by pid
func QueryPipelineVersionByPid(pid int64) (*PipelineVersion, error) {
	var pipelineVersion PipelineVersion
	engineDb := Db
	err := engineDb.Where("pid = ?", pid).Where("status = ?", DataValidStatus).First(&pipelineVersion).Error
	return &pipelineVersion, err
}
