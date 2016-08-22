package etcd

import (
	"testing"
	"time"

	"github.com/containerops/vessel/models"
)

func init() {
	clientEtcd()
}

func TestSaveStage(t *testing.T) {
	stage := fulStage()
	if err := SaveStage(stage); err != nil {
		t.Errorf("Save stage to etcd err : %v", err.Error())
	}
}

func TestGetStage(t *testing.T) {
	stage := easyStage()
	if err := GetStage(stage); err != nil {
		t.Errorf("Get stage from etcd err : %v", err.Error())
	} else {
		t.Log(stage)
	}
}

func TestSetStageStatus(t *testing.T) {
	stage := easyStage()
	stage.Status = models.StateDeleted
	if err := SetStageStatus(stage); err != nil {
		t.Errorf("Save stage status to etcd err : %v", err.Error())
	}
}

func TestGetStageStatus(t *testing.T) {
	stage := easyStage()
	str, err := GetStageStatus(stage)
	if err != nil {
		t.Errorf("Get stage status from etcd err : %v", err.Error())
	} else {
		t.Logf("Stage status in etcd is %v", str)
	}
}

func TestSetStageResult(t *testing.T) {
	result := &models.StageResult{
		Namespace: "etcdStageResult",
		ID:        "bbbbbbbbbb",
		Name:      "stageNamea",
		Result:    models.ResultSuccess,
		Detail:    "VVVVVVVV",
	}
	if err := SetStageResult(result); err != nil {
		t.Errorf("Save stage result to etcd err : %v", err.Error())
	}
}

func TestGetStageResult(t *testing.T) {
	result := &models.StageResult{
		Namespace: "etcdStageResult",
		Name:      "stageNamea",
	}
	if err := GetStageResult(result); err != nil {
		t.Errorf("Get stage result to etcd err : %v", err.Error())
	} else {
		t.Log(result)
	}
}

func TestSetStageTTL(t *testing.T) {
	stage := easyStage()
	if err := SetStageTTL(stage, 2); err != nil {
		t.Errorf("Set stage TTL to etcd err : %v", err.Error())
	} else {
		<-time.After(time.Second * time.Duration(3))
		if err := GetStage(stage); err == nil {
			t.Errorf("Set stage TTL to etcd failed")
		}
	}

}

func easyStage() *models.Stage {
	return &models.Stage{
		Name:      "etcdStage",
		Namespace: "etcdStageResult",
	}
}

func fulStage() *models.Stage {
	return &models.Stage{
		Name:                "etcdStage",
		Namespace:           "etcdStageResult",
		Replicas:            3,
		Image:               "unknow",
		Port:                80,
		StatusCheckURL:      "/heath",
		StatusCheckInterval: 30,
		StatusCheckCount:    3,
		EnvName:             "",
		EnvValue:            "",
		Dependence:          "stageA,stageB,stageC",
		Status:              models.StateRunning,
	}
}
