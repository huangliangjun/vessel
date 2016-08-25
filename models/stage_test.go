package models

import (
	"testing"
)

func Test_AddStageVersion(t *testing.T) {
	stageVersion := &StageVersion{
		Sid:           1,
		Pvid:          1,
		Detail:        "",
		VersionStatus: "ready",
	}
	id, err := AddStageVersion(stageVersion)
	if id == -1 || id == 0 || err != nil {
		t.Error("AddStageVersion failure ", err)
	} else {
		t.Log("AddStageVersion success ", stageVersion)
	}
}

func Test_UpdateStageVersion(t *testing.T) {
	stageVersion := &StageVersion{
		Sid:           1,
		Detail:        "",
		VersionStatus: "running",
	}
	err := UpdateStageVersion(stageVersion)
	if err != nil {
		t.Error("UpdateStageVersion failure ", err)
	} else {
		t.Log("UpdateStageVersion success")
	}
}

func Test_QueryStageVersionBySid(t *testing.T) {
	var sid int64 = 1
	stageVersion, err := QueryStageVersionBySid(sid)
	if err != nil {
		t.Error("GetStageVersionBySid failure ", err)
	} else {
		t.Log("GetStageVersionBySid success ", stageVersion)
	}
}
