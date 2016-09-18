package models

import (
	"time"

	"encoding/json"
	//"fmt"
	"github.com/containerops/vessel/db"
	"github.com/containerops/vessel/utils/timer"
	//	"k8s.io/kubernetes/pkg/api/v1"
)

//Stage Type
const (
	STAGECONTAINER = "container"
	STAGEVM        = "vm"
	STAGEPC        = "pc"
)

// Stage stage data
type Stage struct {
	ID            uint64           `json:"id" gorm:"primary_key"`
	PID           uint64           `json:"pid" gorm:"type:int;not null;index"`
	Namespace     string           `json:"namespace" binding:"Required" gorm:"type:varchar(20);not null;unique_index:idxs_namespace_name"`
	Name          string           `json:"name" binding:"Required" gorm:"type:varchar(20);not null;unique_index:idxs_namespace_name"`
	Replicas      uint             `json:"replicas" gorm:"type:int;default:0"`
	PipelineName  string           `json:"-" sql:"-"`
	Type          string           `json:"type"  binding:"In(container,vm,pc)" gorm:"varchar(20)"`
	Dependencies  string           `json:"dependencies,omitempty" gorm:"varchar(255)"` // eg : "a,b,c"
	Artifacts     []Artifact       `json:"artifacts"  sql:"-"`
	Volumes       []Volume         `json:"volumes"  sql:"-"`
	ArtifactsJSON string           `json:"-" gorm:"column:artifactsJson;type:text;not null"` // json type
	VolumesJSON   string           `json:"-" gorm:"column:volumesJson;type:text;not null"`   // json type
	Ports         []ServicePort    `json:"ports,omitempty" sql:"-"`
	PortsJSON     string           `json:"-" gorm:"column:portsJSON;type:text;not null"` // json type
	Hourglass     *timer.Hourglass `json:"-" sql:"-"`
	Status        uint             `json:"status" gorm:"type:int;default:0"`
	CreatedAt     *time.Time       `json:"created" `
	UpdatedAt     *time.Time       `json:"updated"`
	DeletedAt     *time.Time       `json:"deleted"`
}

// StageVersion data
type StageVersion struct {
	ID           uint64        `json:"id" gorm:"primary_key"`
	PvID         uint64        `json:"pvid" gorm:"type:int;not null"`
	SID          uint64        `json:"sid" gorm:"type:int;not null;index"`
	State        string        `json:"state" gorm:"column:state;type:varchar(20)"`
	Detail       string        `json:"detail" gorm:"type:text;"`
	MetaData     *Stage        `json:"-" sql:"-"`
	PointVersion *PointVersion `json:"-" sql:"-"`
	Status       uint          `json:"status" gorm:"type:tinyint;default:0"`
	CreatedAt    *time.Time    `json:"created" `
	UpdatedAt    *time.Time    `json:"updated"`
	DeletedAt    *time.Time    `json:"deleted"`
}

// Artifact data
type Artifact struct {
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Lifecycle *Lifecycle `json:"lifecycle,omitempty"`
	Container *Container `json:"container,omitempty"`
}

// Lifecycle data
type Lifecycle struct {
	Before  []string `json:"before,omitempty"`
	Runtime []string `json:"runtime,omitempty"`
	After   []string `json:"after,omitempty"`
}

// Container data
type Container struct {
	WorkingDir string          `json:"workingDir,omitempty"`
	Ports      []ContainerPort `json:"ports,omitempty"`
	Env        []EnvVar        `json:"env,omitempty"`
}

// ContainerPort data
type ContainerPort struct {
	Name          string `json:"name,omitempty"`
	HostPort      int32  `json:"hostPort,omitempty"`
	ContainerPort int32  `json:"containerPort,omitempty"`
}

//ServicePort data
type ServicePort struct {
	Name string `json:"name,omitempty"`
	Port int32  `json:"port,omitempty"`
}

// EnvVar data
type EnvVar struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Volume data
type Volume struct {
	Name     string `json:"name,omitempty"`
	HostPath string `json:"hostPath,omitempty"`
}

// StageResult stage result
type StageResult struct {
	Namespace string `json:"-"`
	ID        uint64 `json:"stageVersionID"`
	Name      string `json:"stageName"`
	Status    string `json:"runResult"`
	Detail    string `json:"detail"`
}

// custom set Stage's table name is stage in db
func (s *Stage) TableName() string {
	return "stage"
}

// custom set StageVersion's table name is pipeline_stage_version in db
func (StageVersion) TableName() string {
	return "pipeline_stage_version"
}

//query pipeline'statge data
func (s *Stage) QueryM() ([]*Stage, error) {
	stages := make([]*Stage, 0, 10)
	err := db.Instance.QueryM(s, &stages)
	if err != nil {
		return nil, err
	}
	for i, v := range stages {
		err := v.JsonToObj()
		stages[i] = v
		if err != nil {
			return nil, err
		}
	}
	return stages, err
}

// Obj Artifacts and Volumes To Json
func (s *Stage) ObjToJson() error {
	bsArtifacts, err := json.Marshal(s.Artifacts)
	if err != nil {

		return err
	}
	s.ArtifactsJSON = string(bsArtifacts)

	bsVolumes, err := json.Marshal(s.Volumes)
	if err != nil {
		return err
	}
	s.VolumesJSON = string(bsVolumes)

	bsPorts, err := json.Marshal(s.Ports)
	if err != nil {
		return err
	}
	s.PortsJSON = string(bsPorts)
	return nil
}

// ArtifactsJson and VolumesJson To Obj
func (s *Stage) JsonToObj() error {
	var artifacts []Artifact
	var volumes []Volume
	var ports []ServicePort

	err := json.Unmarshal([]byte(s.ArtifactsJSON), &artifacts)
	if err != nil {
		return err
	}
	s.Artifacts = artifacts
	err = json.Unmarshal([]byte(s.VolumesJSON), &volumes)
	if err != nil {
		return err
	}
	s.Volumes = volumes

	err = json.Unmarshal([]byte(s.PortsJSON), &ports)
	if err != nil {
		return err
	}
	s.Ports = ports
	return nil
}

//create stage version data
func (sv *StageVersion) Create() error {
	if err := db.Instance.Create(sv); err != nil {
		return err
	}
	return nil
}

//update stage version data
func (sv *StageVersion) Update() error {
	if err := db.Instance.Update(sv); err != nil {
		return err
	}
	return nil
}

//query stage version data
func (sv *StageVersion) QueryM() ([]*StageVersion, error) {
	svs := make([]*StageVersion, 0, 10)
	if err := db.Instance.QueryM(sv, &svs); err != nil {
		return nil, err
	}
	return svs, nil
}

//delete stage version data
func (sv *StageVersion) SoftDelete() error {
	if err := db.Instance.DeleteS(sv); err != nil {
		return err
	}
	return nil
}
