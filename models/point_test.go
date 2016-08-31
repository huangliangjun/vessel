package models

import (
	"testing"
)

func Test_PointVersionAdd(t *testing.T) {
	pv := &PointVersion{
		PointID: 1,
		PvID:    1,
		Detail:  "",
		State:   "ready",
	}

	if err := pv.Add(); err != nil {
		t.Error("PointVersion Add failure ", err)
	} else {
		t.Log("PointVersion Add success ", pv)
	}
}

func Test_PointVersionUpdate(t *testing.T) {
	pv := &PointVersion{
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

func Test_PointVersionQueryOne(t *testing.T) {
	pv := &PointVersion{
		PointID: 1,
		Detail:  "",
		State:   "running",
	}

	if err := pv.QueryOne(); err != nil {
		t.Error("PointVersion QueryOne failure ", err)
	} else {
		t.Log("PointVersion QueryOne success ", pv)
	}
}
