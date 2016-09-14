package models

import (
	"testing"
)

func Test_PointVersionCreate(t *testing.T) {
	pv := &PointVersion{
		PointID: 1,
		PvID:    1,
		Detail:  "",
		State:   "ready",
	}

	if err := pv.Create(); err != nil {
		t.Error("PointVersion Create failure ", err)
	} else {
		t.Log("PointVersion Create success ", pv)
	}
}

func Test_PointVersionUpdate(t *testing.T) {
	pv := &PointVersion{
		ID:      1,
		PointID: 1,
		Detail:  "",
		State:   "running",
	}
	if err := pv.Update(); err != nil {
		t.Error("PointVersion Update failure ", err)
	} else {
		t.Log("PointVersion Update success ", pv)
	}
}

func Test_PointVersionQueryM(t *testing.T) {
	pv := &PointVersion{
		PointID: 1,
		State:   "running",
	}
	if pvs, err := pv.QueryM(); err != nil {
		t.Error("PointVersion QueryM failure ", err)
	} else {
		t.Log("PointVersion QueryM success ", pvs)
	}
}
func Test_PointVersionDelete(t *testing.T) {
	pv := &PointVersion{
		PointID: 1,
	}
	if err := pv.SoftDelete(); err != nil {
		t.Error("PointVersion Delete failure ", err)
	} else {
		t.Log("PointVersion Delete success ")
	}
}
