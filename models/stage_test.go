package models

import (
	"testing"
)

func Test_StageVersionCreate(t *testing.T) {
	sv := &StageVersion{
		SID:    1,
		PvID:   1,
		Detail: "",
		State:  "ready",
	}

	if err := sv.Create(); err != nil {
		t.Error("StageVersion Create failure ", err)
	} else {
		t.Log("AddStageVersion Create success ", sv)
	}
}

func Test_StageVersionUpdate(t *testing.T) {
	sv := &StageVersion{
		ID:     1,
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

func Test_StageVersionQueryM(t *testing.T) {
	sv := &StageVersion{
		SID: 1,
	}
	if svs, err := sv.QueryM(); err != nil {
		t.Error("StageVersion QueryM failure ", err)
	} else {
		t.Log("StageVersion QueryM success ", svs)
	}
}

func Test_StageVersionDelete(t *testing.T) {
	sv := &StageVersion{
		SID: 1,
	}
	if err := sv.SoftDelete(); err != nil {
		t.Error("StageVersion Delete failure ", err)
	} else {
		t.Log("StageVersion Delete success ")
	}
}
