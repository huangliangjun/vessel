package models

import (
	"time"

	"encoding/json"
	//	"fmt"

	"k8s.io/kubernetes/pkg/api/v1"
)

const (
	StageContainer StageKind = "container"
	StageVM        StageKind = "vm"
	StagePC        StageKind = "pc"
)

type StageKind string

// Stage stage data
type Stage struct {
	ID            int64       `json:"id" gorm:"primary_key"`
	Pid           int64       `json:"pid" gorm:"type:bigint;not null;index"`
	Name          string      `json:"name"  binding:"Required" gorm:"type:varchar(30);unique"`
	Type          string      `json:"type"  binding:"In(container,vm,pc)" gorm:"type:varchar(20)"`
	Dependencies  string      `json:"dependencies,omitempty" gorm:"varchar(255)"` // eg : "a,b,c"
	Artifacts     []*Artifact `json:"artifacts" binding:"Required" gorm:"-"`
	Volumes       []*Volume   `json:"volumes,omitempty" gorm:"-"`
	ArtifactsJson string      `json:"-" gorm:"column:artifactsJson;type:text;not null"` // json type
	VolumesJson   string      `json:"-" gorm:"column:volumesJson;type:text;not null"`   // json type
	Status        uint        `json:"status" gorm:"type:tinyint;default:0"`
	CreatedAt     *time.Time  `json:"created" `
	UpdatedAt     *time.Time  `json:"updated"`
	DeletedAt     *time.Time  `json:"deleted"`
}

// StageVersion data
type StageVersion struct {
	Id            int64      `json:"id" gorm:"primary_key"`
	Sid           int64      `json:"sid" gorm:"type:bigint;not null"`
	Pvid          int64      `json:"pvid" gorm:"type:bigint;not null"`
	Detail        string     `json:"detail" gorm:"type:text;"`
	VersionStatus string     `json:"versionStatus" gorm:"column:versionStatus;type:varchar(20);not null;"`
	Status        uint       `json:"status" gorm:"type:tinyint;default:0"`
	CreatedAt     *time.Time `json:"created" `
	UpdatedAt     *time.Time `json:"updated"`
	DeletedAt     *time.Time `json:"deleted"`
}

type Artifact struct {
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Lifecycle *Lifecycle `json:"lifecycle,omitempty"`
	Container *Container `json:"container,omitempty"`
}

type Lifecycle struct {
	Before  []string `json:"before,omitempty"`
	Runtime []string `json:"runtime,omitempty"`
	After   []string `json:"after,omitempty"`
}

type Volume struct {
	Name                 string                               `json:"name"`
	HostPath             *v1.HostPathVolumeSource             `json:"hostPath,omitempty"`
	EmptyDir             *v1.EmptyDirVolumeSource             `json:"emptyDir,omitempty"`
	AWSElasticBlockStore *v1.AWSElasticBlockStoreVolumeSource `json:"awsElasticBlockStore,omitempty"`
	CephFS               *v1.CephFSVolumeSource               `json:"cephfs,omitempty"`
}

type Container struct {
	WorkingDir     string             `json:"workingDir,omitempty"`
	Ports          []v1.ContainerPort `json:"ports,omitempty"`
	Env            []v1.EnvVar        `json:"env,omitempty"`
	VolumeMounts   []v1.VolumeMount   `json:"volumeMounts,omitempty"`
	LivenessProbe  *v1.Probe          `json:"livenessProbe,omitempty"`
	ReadinessProbe *v1.Probe          `json:"readinessProbe,omitempty"`
	PullPolicy     v1.PullPolicy      `json:"PullPolicy,omitempty"`
	Stdin          bool               `json:"stdin,omitempty"`
	StdinOnce      bool               `json:"stdinOnce,omitempty"`
	TTY            bool               `json:"tty,omitempty"`
}

//custom set Stage's table name to be pipeline_stage
func (Stage) TableName() string {
	return "pipeline_stage"
}

//custom set StageVersion's table name to be pipeline_stage_version
func (StageVersion) TableName() string {
	return "pipeline_stage_version"
}

//query pipeline'statge data
func (s *Stage) Query() ([]*Stage, error) {
	engineDb := Db
	stages := make([]*Stage, 0, 10)
	s.Status = DataValidStatus
	err := engineDb.Find(&stages, s).Error
	return stages, err
}

// Obj Artifacts and Volumes To Json
func (s *Stage) ObjToJson() error {
	bsArtifacts, err := json.Marshal(s.Artifacts)
	if err != nil {

		return err
	}
	s.ArtifactsJson = string(bsArtifacts)

	bsVolumes, err := json.Marshal(s.Volumes)
	if err != nil {
		return err
	}
	s.VolumesJson = string(bsVolumes)
	return nil
}

// ArtifactsJson and VolumesJson To Obj
func (s *Stage) JsonToObj() error {
	var artifacts []*Artifact
	var volumes []*Volume
	err := json.Unmarshal([]byte(s.ArtifactsJson), &artifacts)
	if err != nil {
		return err
	}
	s.Artifacts = artifacts
	err = json.Unmarshal([]byte(s.VolumesJson), &volumes)
	if err != nil {
		return err
	}
	s.Volumes = volumes
	return nil
}

//save stage version data
func (sv *StageVersion) Add() error {
	engineDb := Db
	return engineDb.Create(sv).Error
}

//update stage version data
func (sv *StageVersion) Update() error {
	engineDb := Db
	return engineDb.Model(sv).Update(sv).Error
}

//query stage version data
func (sv *StageVersion) QueryOne() error {
	engineDb := Db
	return engineDb.First(sv).Error
}
