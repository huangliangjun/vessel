package models

import (
	"testing"
)

func Test_StageVersionAdd(t *testing.T) {
	sv := &StageVersion{
		Sid:           1,
		Pvid:          1,
		Detail:        "",
		VersionStatus: "ready",
	}

	if err := sv.Add(); err != nil {
		t.Error("StageVersion Add failure ", err)
	} else {
		t.Log("AddStageVersion success ", sv)
	}
}

func Test_StageVersionUpdate(t *testing.T) {
	sv := &StageVersion{
		Sid:           1,
		Detail:        "",
		VersionStatus: "running",
	}

	if err := sv.Update(); err != nil {
		t.Error("StageVersion Update failure ", err)
	} else {
		t.Log("StageVersion Update success", sv)
	}
}

func Test_StageVersionQueryOne(t *testing.T) {
	sv := &StageVersion{
		Sid:           1,
		Detail:        "",
		VersionStatus: "running",
	}

	if err := sv.QueryOne(); err != nil {
		t.Error("StageVersion QueryOne failure ", err)
	} else {
		t.Log("StageVersion QueryOne success ", sv)
	}
}
