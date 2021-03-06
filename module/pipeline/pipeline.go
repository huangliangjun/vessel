package pipeline

import (
	"encoding/json"
	"log"

	//	"errors"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/dependence"
	"github.com/containerops/vessel/module/point"
	"github.com/containerops/vessel/module/scheduler"
)

// CreatePipeline new pipeline with PipelineTemplate
func CreatePipeline(pipelineTemplate *models.PipelineTemplate) []byte {
	log.Println("Create pipeline")
	pipeline := pipelineTemplate.MetaData

	// Check pipeline exist
	if is := pipeline.CheckIsExist(); is == true {
		bytes, _ := outputResult(pipeline, 0, nil, "pipeline had exist.")
		return bytes
	}
	log.Println("the pipeline is ", pipeline)

	// Check stages and points dependence
	if err := dependence.CheckDependence(pipeline); err != nil {
		bytes, _ := outputResult(pipeline, 0, nil, err.Error())
		return bytes
	}

	// Insert pipeline
	if err := pipeline.Add(); err != nil {
		bytes, _ := outputResult(pipeline, 0, nil, err.Error())
		return bytes
	}

	log.Printf("Create pipeline name = %v in namespace '%v' is over", pipeline.Namespace, pipeline.Name)
	log.Print("Create job is done")
	bytes, _ := outputResult(pipeline, 0, nil, "")
	return bytes

}

//StartPipeline start pipeline in k8s
func StartPipeline(pID uint64) []byte {
	log.Println("Start pipeline")
	pipeline := &models.Pipeline{ID: pID}

	// get pipeline data
	if err := pipeline.QueryOne(); err != nil {
		bytes, _ := outputResult(pipeline, 0, nil, err.Error())
		return bytes
	}

	pipelineVsn := &models.PipelineVersion{
		PID:      pipeline.ID,
		State:    models.StateReady,
		MetaData: pipeline,
		Status:   models.DataValidStatus,
	}

	// Insert pipelineVersion
	if err := pipelineVsn.Add(); err != nil {
		bytes, _ := outputResult(pipeline, 0, nil, err.Error())
		return bytes
	}
	log.Println(pipelineVsn)
	executorMap := dependence.ParsePipelineVersion(pipelineVsn)
	bs, _ := json.MarshalIndent(executorMap, " ", "  ")
	log.Println(string(bs))
	schedulingRes := scheduler.StartPoint(executorMap, models.StartPointMark)
	bytes, success := outputResult(pipeline, pipelineVsn.ID, schedulingRes, "")
	if success {
		pipelineVsn.State = models.StateRunning
		// TODO:update pipeline version status
		if err := pipelineVsn.Update(); err != nil {
			bytes, _ := outputResult(pipeline, 0, nil, err.Error())
			return bytes
		}

	} else {
		//rollback by pipeline failed
		go removePipeline(executorMap, pipelineVsn, "run pipeline error")
	}

	byteStr, _ := json.Marshal(pipeline)
	log.Println(string(byteStr))
	log.Printf("Start pipeline name = %v in namespace '%v' is over", pipeline.Namespace, pipeline.Name)
	log.Print("Start job is done")
	return bytes
	//return nil
}

//StopPipeline stop pipeline in k8s
func StopPipeline(pID uint64, pvID uint64) []byte {
	log.Println("Stop pipeline")
	// TODO: Get pipeline form db
	pipeline := &models.Pipeline{
		ID: pID,
	}
	if err := pipeline.QueryOne(); err != nil {
		bytes, _ := outputResult(pipeline, 0, nil, err.Error())
		return bytes
	}
	// TODO: Get pipeline version form db
	pipelineVersion := &models.PipelineVersion{
		ID:       pvID,
		MetaData: pipeline,
	}
	if err := pipelineVersion.QueryOne(); err != nil {
		bytes, _ := outputResult(pipeline, 0, nil, err.Error())
		return bytes
	}
	executorMap := dependence.ParsePipelineVersion(pipelineVersion)

	//point version delete
	if err := point.Delete(pvID); err != nil {
		bytes, _ := outputResult(pipeline, 0, nil, err.Error())
		return bytes
	}

	//stop pipeline
	schedulingRes := removePipeline(executorMap, pipelineVersion, "")
	bytes, _ := outputResult(pipeline, 0, schedulingRes, "")
	log.Printf("Delete pipeline name = %v in namespace '%v' is over", pipeline.Namespace, pipeline.Name)
	log.Print("Delete job is done")
	return bytes
}

// DeletePipeline delete pipeline by pID
func DeletePipeline(pID uint64) []byte {
	log.Println("Delete pipeline")
	// TODO: Get pipeline form db
	pipeline := &models.Pipeline{
		ID:     pID,
		Status: models.DataValidStatus,
	}
	// Check pipeline exist
	if is := pipeline.CheckIsExist(); is == false {
		bytes, _ := outputResult(pipeline, 0, nil, "pipeline not exist")
		return bytes
	}
	// TODO: Get pipeline version list form db with pID when is not delete
	pipelineVsn := &models.PipelineVersion{
		PID:    pID,
		Status: models.DataValidStatus,
	}
	if err := pipelineVsn.QueryOne(); err == nil {
		bytes, _ := outputResult(pipeline, 0, nil, "pipeline is running,could not delete")
		return bytes
	}

	// TODO:delete pipeline
	if err := pipeline.Delete(); err != nil {
		bytes, _ := outputResult(pipeline, 0, nil, "pipeline delete failure")
		return bytes
	}
	log.Printf("Delete pipeline name = %v in namespace '%v' is over", pipeline.Namespace, pipeline.Name)
	log.Print("Delete job is done")
	bytes, _ := outputResult(pipeline, 0, nil, "")
	return bytes
}

func RenewPipeline(pID uint64, pipelineTemplate *models.PipelineTemplate) []byte {
	log.Println("Renew pipeline")
	pipeline := pipelineTemplate.MetaData

	// TODO:check pipeline already exist
	//if namespace name pipeline not in db {
	//	detail := fmt.Sprintf("Pipeline = %v in namespane = %v is not already exist", pipeline.Name, pipeline.Namespace)
	//	bytes, _ := formatOutputBytes(pipelineTemplate, pipeline, nil, detail)
	//	return bytes
	//}

	// TODO:check stage already exist
	//for _, stage := range pipeline.Stages {
	//if namespace name stage not in db {
	//	detail := fmt.Sprintf("Stage = %v in namespane = %v is not already exist", stage.Name, stage.Namespace)
	//	bytes, _ := formatOutputBytes(pipelineTemplate, pipeline, nil, detail)
	//	return bytes
	//}
	//}

	//if err := dependence.CheckPipeline(pipeline); err != nil {
	//	bytes, _ := outputResult(pipeline, 0, nil, err.Error())
	//	return bytes
	//}

	// TODO:update all pipeline with pID
	log.Printf("Renew pipeline name = %v in namespace '%v' is over", pipeline.Namespace, pipeline.Name)
	log.Print("Renew job is done")
	bytes, _ := outputResult(pipeline, 0, nil, "")
	return bytes
	return nil
}

//get pipeline info by pid
func GetPipeline(pID uint64) []byte {
	log.Println("Renew pipeline")
	// TODO:Get pipeline
	return nil
}

func removePipeline(executorList []interface{}, pipelineVersion *models.PipelineVersion, detail string) []*models.ExecutedResult {
	schedulingRes := scheduler.StopPoint(executorList, models.StartPointMark)
	// TODO: delete pipeline version
	pipelineVersion.State = models.StateDeleted
	pipelineVersion.Detail = detail
	pipelineVersion.Status = models.DataInValidStatus
	if err := pipelineVersion.Update(); err != nil {
		return schedulingRes
	}
	return schedulingRes
}

//format output result
func outputResult(pipeline *models.Pipeline, pvID uint64, schedulingRes []*models.ExecutedResult, detail string) ([]byte, bool) {
	log.Println("Pipeline result :", schedulingRes)
	status := models.ResultFailed
	if detail == "" {
		status = models.ResultSuccess
		if schedulingRes != nil {
			for _, result := range schedulingRes {
				if status != result.Status {
					status = result.Status
					detail = result.Detail
					break
				}
			}
		}
	}
	output := &models.PipelineResult{
		PID:       pipeline.ID,
		Name:      pipeline.Name,
		Namespace: pipeline.Namespace,
		Status:    status,
		Detail:    detail,
	}
	if pvID != 0 {
		output.PvID = pvID
	}
	bytes, err := json.Marshal(output)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Pipeline result is %v", string(bytes))
	return bytes, status == models.ResultSuccess
}
