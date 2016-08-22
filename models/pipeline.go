package models

import "time"

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
	Stages     []*Stage     `json:"stages"`
	Points     []*Point     `json:"points"`
}

// Pipeline pipeline data
type Pipeline struct {
	ID        int64      `json:"id" gorm:"primary_key"`
	Namespace string     `json:"namespace" binding:"Required"`
	Name      string     `json:"name" binding:"Required"`
	Timeout   uint64     `json:"timeout"`
	Status    uint       `json:"Status" gorm:`
	Created   *time.Time `json:"created" `
	Updated   *time.Time `json:"updated"`
	Deleted   *time.Time `json:"deleted"`
}

// PipelineVersion data
type PipelineVersion struct {
	ID      int64      `json:"id" gorm:"primary_key"`
	Pid     int64      `json:"Pid" gorm:`
	Status  uint       `json:"status" gorm:`
	Detail  string     `json:"detail" gorm:`
	Created *time.Time `json:"created" `
	Updated *time.Time `json:"updated"`
	Deleted *time.Time `json:"deleted"`
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
