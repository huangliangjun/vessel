package point

import (
	//"fmt"

	"github.com/containerops/vessel/models"
)

func StartPoint(pointVsn *models.PointVersion, readyMap map[string]bool) (bool, bool) {
	meet := true
	ended := false
	if pointVsn.Kind == models.StartPoint {
		err := AddAndUpdate(pointVsn)
		if err != nil {
			return meet, ended
		}
		return meet, ended
	}
	for _, condition := range pointVsn.Conditions {
		if meet := readyMap[condition]; !meet {
			return meet, ended
		}
	}
	ended = pointVsn.Kind == models.EndPoint
	err := AddAndUpdate(pointVsn)
	if err != nil {
		return meet, ended
	}

	return meet, ended
}

func StopPoint(pointVsn *models.PointVersion, readyMap map[string]bool) (bool, bool) {
	meet := true
	ended := false
	if pointVsn.Kind == models.StartPoint {
		err := Delete(pointVsn)
		if err != nil {
			return meet, ended
		}
		return meet, ended
	}
	for _, condition := range pointVsn.Conditions {
		if meet := readyMap[condition]; !meet {
			return meet, ended
		}
	}
	ended = pointVsn.Kind == models.EndPoint
	err := Delete(pointVsn)
	if err != nil {
		return meet, ended
	}
	return meet, ended
}

func AddAndUpdate(pointVsn *models.PointVersion) error {
	if pointVsn.Kind != models.TemporaryPoint {
		pointVsn.State = models.StateReady
		if err := pointVsn.Create(); err != nil {
			return err
		}
		pointVsn.State = models.StateRunning
		if err := pointVsn.Update(); err != nil {
			return err
		}
	}

	return nil
}

func Delete(pointVsn *models.PointVersion) error {
	pv := &models.PointVersion{
		PvID:    pointVsn.PvID,
		PointID: pointVsn.PointID,
		State:   models.StateDeleted,
	}
	if err := pv.Update(); err != nil {
		return err
	}
	pv = &models.PointVersion{
		PvID:    pointVsn.PvID,
		PointID: pointVsn.PointID,
	}
	if err := pv.SoftDelete(); err != nil {
		return err
	}
	return nil
}
