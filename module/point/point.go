package point

import (
	//	"fmt"

	"github.com/containerops/vessel/models"
)

func CheckPoint(pointVsn *models.PointVersion, readyMap map[string]bool) (bool, bool) {
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
	err := AddAndUpdate(pointVsn)
	if err != nil {
		return meet, ended
	}
	ended = pointVsn.Kind == models.EndPoint
	return meet, ended
}

func AddAndUpdate(pointVsn *models.PointVersion) error {
	if pointVsn.Kind != models.TemporaryPoint {
		pointVsn.State = models.StateReady
		pointVsn.Status = models.DataValidStatus
		if err := pointVsn.Add(); err != nil {
			return err
		}
		pointVsn.State = models.StateRunning
		if err := pointVsn.Update(); err != nil {
			return err
		}
	}

	return nil
}

func Delete(pvid uint64) error {
	pointVsn := &models.PointVersion{
		PvID:   pvid,
		State:  models.StateDeleted,
		Status: models.DataInValidStatus,
	}
	return pointVsn.Update()
}
