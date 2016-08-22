package stage

import (
	"fmt"
	"log"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/etcd"
	kubeclt "github.com/containerops/vessel/module/kubernetes"
	"github.com/containerops/vessel/utils"
	"github.com/containerops/vessel/utils/timer"
)

// StartStage start stage workflow
func StartStage(info interface{}, finishChan chan *models.ExecutedResult, hourglass *timer.Hourglass) {
	stage := info.(*models.Stage)
	if stage == nil {
		finishChan <- fillSchedulingResult(stage, models.ResultFailed, "Stage is empty")
		return
	}

	etcd.SaveStage(stage)
	stage.Status = models.StateReady
	etcd.SetStageStatus(stage)

	res := kubeclt.CreateStage(stage, hourglass)
	if res.Result != models.ResultSuccess {
		finishChan <- fillSchedulingResult(stage, res.Result, res.Detail)
		return
	}
	res = kubeclt.GetBusinessRes(stage, hourglass)
	if res.Result != models.ResultSuccess {
		finishChan <- fillSchedulingResult(stage, res.Result, res.Detail)
		return
	}
	stage.Status = models.StateRunning
	etcd.SetStageStatus(stage)
	finishChan <- fillSchedulingResult(stage, models.ResultSuccess, "")
}

// StopStage stop stage workflow
func StopStage(info interface{}, finishChan chan *models.ExecutedResult, hourglass *timer.Hourglass) {
	stage := info.(*models.Stage)
	if stage == nil {
		finishChan <- fillSchedulingResult(stage, models.ResultFailed, "Stage is empty")
		return
	}

	res := kubeclt.DeleteStage(stage, hourglass)
	stage.Status = models.StateDeleted
	etcd.SetStageStatus(stage)
	finishChan <- fillSchedulingResult(stage, res.Result, res.Detail)
}

func fillSchedulingResult(stage *models.Stage, result string, detail string) *models.ExecutedResult {
	log.Println(fmt.Sprintf("Stage name = %v result is %v, detail is %v", stage.Name, result, detail))
	stageName := ""
	namespace := ""
	if stage != nil {
		stageName = stage.Name
		namespace = stage.Namespace
	}
	return &models.ExecutedResult{
		Name:   stageName,
		Status: result,
		Result: &models.StageResult{
			ID:        utils.UUID(),
			Namespace: namespace,
			Name:      stageName,
			Result:    result,
			Detail:    detail,
		},
	}
}
