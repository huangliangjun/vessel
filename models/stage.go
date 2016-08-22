package models

import "k8s.io/kubernetes/pkg/api/v1"

const (
	StageContainer StageKind = "container"
	StageVM        StageKind = "vm"
	StagePC        StageKind = "pc"
)

type StageKind string

// Stage stage data
type Stage struct {
	ID            int64       `json:"-"`
	Name          string      `json:"name"  binding:"Required"`
	Type          StageKind   `json:"type"  binding:"In(container,vm,pc)"`
	Dependencies  string      `json:"dependencies,omitempty"` // eg : "a,b,c"
	Artifacts     []*Artifact `json:"artifacts" binding:"Required"`
	Volumes       []*Volume   `json:"volumes,omitempty"`
	ArtifactsJson string      `json:"-" ` // json type
	VolumesJson   string      `json:"-"`  // json type
}

// StageVersion data
type StageVersion struct {
	ID      int64      `json:"id" gorm:"primary_key"`
	Sid     int64      `json:"sid" gorm:`
	Pvid    int64      `json:"pvid" gorm:`
	Status  uint       `json:"status" gorm:`
	Detail  string     `json:"detail" gorm:`
	Created *time.Time `json:"created" `
	Updated *time.Time `json:"updated"`
	Deleted *time.Time `json:"deleted"`
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
