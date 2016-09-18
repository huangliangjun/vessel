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
		return meet, ended
	}
	for _, condition := range pointVsn.Conditions {
		if meet := readyMap[condition]; !meet {
			return meet, ended
		}
	}
	ended = pointVsn.Kind == models.EndPoint
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

func Delete(pvid uint64) error {
	pointVsn := &models.PointVersion{
		PvID:  pvid,
		State: models.StateDeleted,
	}
	if err := pointVsn.Update(); err != nil {
		return err
	}
	pointVsn = &models.PointVersion{
		PvID: pvid,
	}
	if err := pointVsn.SoftDelete(); err != nil {
		return err
	}
	return nil
}
