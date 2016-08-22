package etcd

import (
	"log"
	"testing"
	"time"

	"github.com/containerops/vessel/models"
)

func init() {
	clientEtcd()
}

func TestSavePipeline(t *testing.T) {
	pipeline := fulPipeline()
	if err := SavePipeline(pipeline); err != nil {
		t.Errorf("Save pipeline to etcd err : %v", err.Error())
	}
}

func TestGetPipeline(t *testing.T) {
	pipeline := easyPipeline()
	if err := GetPipeline(pipeline); err != nil {
		t.Errorf("Get pipeline to etcd err : %v", err.Error())
	} else {
		t.Log(pipeline)
	}
}

func TestSetPipelineStatus(t *testing.T) {
	pipeline := fulPipeline()
	if err := SetPipelineStatus(pipeline); err != nil {
		t.Errorf("Save pipeline status to etcd err : %v", err.Error())
	}
}

func TestGetPipelineStatus(t *testing.T) {
	pipeline := easyPipeline()
	str, err := GetPipelineStatus(pipeline)
	if err != nil {
		t.Errorf("Get pipeline status to etcd err : %v", err.Error())
	} else {
		t.Logf("Pileline status in etcd is %v", str)
	}
}

func TestSetCreationTimestamp(t *testing.T) {
	pipeline := easyPipeline()
	if err := SetCreationTimestamp(pipeline); err != nil {
		t.Errorf("Save pipeline creation timestamp to etcd err : %v", err.Error())
	}
}

func TestGetCreationTimestamp(t *testing.T) {
	pipeline := easyPipeline()
	str, err := GetCreationTimestamp(pipeline)
	if err != nil {
		t.Errorf("Get pipeline createion timestamp from etcd err : %v", err.Error())
	} else {
		t.Logf("Pileline createion timestamp in etcd is %v", str)
	}
}

func TestSetDeletionTimestamp(t *testing.T) {
	pipeline := easyPipeline()
	log.Println(pipeline, SetDeletionTimestamp(pipeline))
	if err := SetDeletionTimestamp(pipeline); err != nil {
		t.Errorf("Set pipeline deletion timestamp from etcd err : %v", err.Error())
	}
}

func TestGetDeletionTimestamp(t *testing.T) {
	pipeline := easyPipeline()
	str, err := GetDeletionTimestamp(pipeline)
	if err != nil {
		t.Errorf("Get pipeline deletion timestamp from etcd err : %v", err.Error())
	} else {
		t.Logf("Pileline deletion timestamp in etcd is %v", str)
	}
}

func TestSetPipelineResult(t *testing.T) {
	result := &models.PipelineResult{
		Name:        "etcdPipelineResult",
		Namespace:   "xx",
		WorkspaceID: 1000,
		ProjectID:   2000,
		PipelineID:  "aaaaaaaaaaa",
		Status:      models.StateReady,
	}
	if err := SetPipelineResult(result); err != nil {
		t.Errorf("Save pipeline result to etcd err : %v", err.Error())
	}
}

func TestGetPipelineResult(t *testing.T) {
	result := &models.PipelineResult{
		Name:      "etcdPipelineResult",
		Namespace: "xx",
	}
	if err := GetPipelineResult(result); err != nil {
		t.Errorf("Get pipeline result from etcd err : %v", err.Error())
	} else {
		t.Log(result)
	}
}

func TestSetPipelineTTL(t *testing.T) {
	pipeline := easyPipeline()
	if err := SetPipelineTTL(pipeline, 2); err != nil {
		t.Errorf("Set pipeline TTL to etcd err : %v", err.Error())
	} else {
		<-time.After(time.Second * time.Duration(3))
		if err := GetPipeline(pipeline); err == nil {
			t.Errorf("Set pipeline TTL to etcd failed")
		}
	}
}

func easyPipeline() *models.Pipeline {
	return &models.Pipeline{
		Name:      "etcdPipeline",
		Namespace: "xx",
	}
}

func fulPipeline() *models.Pipeline {
	return &models.Pipeline{
		Name:            "etcdPipeline",
		Namespace:       "xx",
		Stages:          []string{"pipelineA", "pipelineB", "pipelineC"},
		TimeoutDuration: 60,
		Status:          models.StateReady,
	}
}
