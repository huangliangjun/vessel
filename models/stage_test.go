package models

import (
	"testing"
)

func Test_StageVersionAdd(t *testing.T) {
	sv := &StageVersion{
		SID:    1,
		PvID:   1,
		Detail: "",
		State:  "ready",
	}

	if err := sv.Add(); err != nil {
		t.Error("StageVersion Add failure ", err)
	} else {
		t.Log("AddStageVersion success ", sv)
	}
}

func Test_StageVersionUpdate(t *testing.T) {
	sv := &StageVersion{
		SID:    1,
		Detail: "",
		State:  "running",
	}

	if err := sv.Update(); err != nil {
		t.Error("StageVersion Update failure ", err)
	} else {
		t.Log("StageVersion Update success", sv)
	}
}

func Test_StageVersionQueryOne(t *testing.T) {
	sv := &StageVersion{
		SID:    1,
		Detail: "",
		State:  "running",
	}

	if err := sv.QueryOne(); err != nil {
		t.Error("StageVersion QueryOne failure ", err)
	} else {
		t.Log("StageVersion QueryOne success ", sv)
	}
}
