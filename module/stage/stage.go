package stage

import (
	"fmt"
	"log"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/deployment"
	"github.com/containerops/vessel/module/point"
)

// Start stage
func Start(info interface{}, readyMap map[string]bool, finishChan chan *models.ExecutedResult) {
	stageVsn := info.(*models.StageVersion)
	metaData := stageVsn.MetaData
	if stageVsn.State != "" {
		return
	}
	meet, ended := point.StartPoint(stageVsn.PointVersion, readyMap)
	if ended {
		log.Println("endPointMark")
		finishChan <- FillSchedulingResult(stageVsn.ID, models.EndPointMark, models.ResultSuccess, "")
		return
	}
	if !meet {
		return
	}
	readyMap[metaData.Name] = true
	//TODO:Save stageVersion
	stageVsn.State = models.StateReady
	stageVsn.Status = models.DataValidStatus
	if err := stageVsn.Add(); err != nil {
		return
	}
	deployment := deployment.NewDeployment(metaData)
	res := deployment.Deploy()
	if res.Status != models.ResultSuccess {
		finishChan <- FillSchedulingResult(stageVsn.ID, metaData.Name, res.Status, res.Detail)
		return
	}

	//TODO:Update stageVersion
	stageVsn.State = models.StateRunning
	if err := stageVsn.Update(); err != nil {
		finishChan <- FillSchedulingResult(stageVsn.ID, metaData.Name, res.Status, res.Detail)
		return
	}
	finishChan <- FillSchedulingResult(stageVsn.ID, stageVsn.MetaData.Name, models.ResultSuccess, "")
}

// Stop stage
func Stop(info interface{}, readyMap map[string]bool, finishChan chan *models.ExecutedResult) {
	stageVsn := info.(*models.StageVersion)
	metaData := stageVsn.MetaData
	if stageVsn.State == models.StateDeleted {
		return
	}
	meet, ended := point.StopPoint(stageVsn.PointVersion, readyMap)
	if ended {
		log.Println("endPointMark")
		finishChan <- FillSchedulingResult(stageVsn.ID, models.EndPointMark, models.ResultSuccess, "")
		return
	}
	if !meet {
		return
	}
	readyMap[metaData.Name] = true
	deployment := deployment.NewDeployment(metaData)
	res := deployment.Undeploy()
	//TODO:Update stageVersion
	stageVsn.State = models.StateDeleted
	stageVsn.Status = models.DataInValidStatus
	if err := stageVsn.Update(); err != nil {
		finishChan <- FillSchedulingResult(stageVsn.ID, stageVsn.MetaData.Name, models.ResultFailed, "stageVersion update failure")
	} else {
		finishChan <- FillSchedulingResult(stageVsn.ID, res.Name, res.Status, res.Detail)
	}

}

func FillSchedulingResult(svid uint64, stageName, result string, detail string) *models.ExecutedResult {
	log.Println(fmt.Sprintf("Stage name = %v result is %v, detail is %v", stageName, result, detail))
	return &models.ExecutedResult{
		SID:    svid,
		Name:   stageName,
		Status: result,
		Detail: detail,
	}
}
