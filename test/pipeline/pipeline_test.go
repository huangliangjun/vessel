package pipeline

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/containerops/vessel/models"
)

//func Test_SingleCreateDeletePipeline(t *testing.T) {
//	path := "./pipeline.json"
//	f, err := os.Open(path)
//	if err != nil {
//		t.Errorf("Read path error : %v", err)
//	}
//	body, code, err := postPipeline(f)
//	if err != nil {
//		t.Errorf("postPipeline error : %v", err)
//	}
//	t.Log("the status code is ", code)
//	t.Log("the response body :", string(body))

//	pipelineResult := models.PipelineResult{}
//	err = json.Unmarshal(body, &pipelineResult)
//	if err != nil {
//		t.Errorf("Unmarshal body error : %v", err)
//	}
//	body, code, err = deletePipeline(pipelineResult.PID)
//	if err != nil {
//		t.Errorf("deletePipeline error : %v", err)
//	}
//	t.Log("the status code is ", code)
//	t.Log("the response body :", string(body))
//}

//func Test_SingleStartStopPipeline(t *testing.T) {
//	path := "./pipeline1.json"
//	f, err := os.Open(path)
//	if err != nil {
//		t.Errorf("Read path error : %v", err)
//	}
//	body, code, err := postPipeline(f)
//	if err != nil {
//		t.Errorf("postPipeline error : %v", err)
//	}
//	pipelineResult := models.PipelineResult{}
//	err = json.Unmarshal(body, &pipelineResult)
//	if err != nil {
//		t.Errorf("Unmarshal body error : %v", err)
//	}
//	body, code, err = startPipeline(pipelineResult.PID)
//	if err != nil {
//		t.Errorf("startPipeline error : %v", err)
//	}
//	t.Log("the status code is ", code)
//	t.Log("the response body :", string(body))
//	err = json.Unmarshal(body, &pipelineResult)
//	if err != nil {
//		t.Errorf("Unmarshal body error : %v", err)
//	}
//	body, code, err = stopPipeline(pipelineResult.PID, pipelineResult.PvID)
//	if err != nil {
//		t.Errorf("stopPipeline error : %v", err)
//	}
//	t.Log("the status code is ", code)
//	t.Log("the response body :", string(body))
//}

//func Test_MultiCreateDeletePipeline(t *testing.T) {
//	path := "./pipeline2.json"
//	f, err := os.Open(path)
//	if err != nil {
//		t.Errorf("Read path error : %v", err)
//	}
//	path1 := "./pipeline3.json"
//	f1, err := os.Open(path1)
//	if err != nil {
//		t.Errorf("Read path error : %v", err)
//	}
//	body, code, err := postPipeline(f)
//	if err != nil {
//		t.Errorf("postPipeline error : %v", err)
//	}
//	t.Log("the status code is ", code)
//	t.Log("the response body :", string(body))
//	body1, code1, err := postPipeline(f1)
//	if err != nil {
//		t.Errorf("postPipeline error : %v", err)
//	}
//	t.Log("the status code is ", code1)
//	t.Log("the response body :", string(body1))

//	pipelineResult := models.PipelineResult{}
//	err = json.Unmarshal(body, &pipelineResult)
//	if err != nil {
//		t.Errorf("Unmarshal body error : %v", err)
//	}
//	pipelineResult1 := models.PipelineResult{}
//	err = json.Unmarshal(body1, &pipelineResult)
//	if err != nil {
//		t.Errorf("Unmarshal body error : %v", err)
//	}
//	body, code, err = deletePipeline(pipelineResult.PID)
//	if err != nil {
//		t.Errorf("deletePipeline error : %v", err)
//	}
//	t.Log("the status code is ", code)
//	t.Log("the response body :", string(body))
//	body1, code1, err = deletePipeline(pipelineResult1.PID)
//	if err != nil {
//		t.Errorf("deletePipeline error : %v", err)
//	}
//	t.Log("the status code is ", code1)
//	t.Log("the response body :", string(body1))
//}

func Test_MultiStartStopPipeline(t *testing.T) {
	path := "./pipeline4.json"
	f, err := os.Open(path)
	if err != nil {
		t.Errorf("Read path error : %v", err)
	}
	path1 := "./pipeline5.json"
	f1, err := os.Open(path1)
	if err != nil {
		t.Errorf("Read path error : %v", err)
	}
	body, code, err := postPipeline(f)
	if err != nil {
		t.Errorf("postPipeline error : %v", err)
	}
	body1, code1, err := postPipeline(f1)
	if err != nil {
		t.Errorf("postPipeline error : %v", err)
	}
	t.Log("the status code is ", code)
	t.Log("the response body :", string(body))
	t.Log("the status code is ", code1)
	t.Log("the response body :", string(body1))
	pipelineResult := models.PipelineResult{}
	err = json.Unmarshal(body, &pipelineResult)
	if err != nil {
		t.Errorf("Unmarshal body error : %v", err)
	}
	t.Log("the pipelineResult is ", pipelineResult)
	pipelineResult1 := models.PipelineResult{}
	err = json.Unmarshal(body1, &pipelineResult1)
	if err != nil {
		t.Errorf("Unmarshal body error : %v", err)
	}
	t.Log("the pipelineResult1 is ", pipelineResult1)
	body, code, err = startPipeline(pipelineResult.PID)
	if err != nil {
		t.Errorf("startPipeline error : %v", err)
	}
	body1, code1, err = startPipeline(pipelineResult1.PID)
	if err != nil {
		t.Errorf("startPipeline error : %v", err)
	}
	t.Log("the status code is ", code)
	t.Log("the response body :", string(body))
	t.Log("the status code is ", code1)
	t.Log("the response body :", string(body1))
	err = json.Unmarshal(body, &pipelineResult)
	if err != nil {
		t.Errorf("Unmarshal body error : %v", err)
	}
	err = json.Unmarshal(body1, &pipelineResult1)
	if err != nil {
		t.Errorf("Unmarshal body error : %v", err)
	}
	body, code, err = stopPipeline(pipelineResult.PID, pipelineResult.PvID)
	if err != nil {
		t.Errorf("stopPipeline error : %v", err)
	}
	body1, code1, err = stopPipeline(pipelineResult1.PID, pipelineResult1.PvID)
	if err != nil {
		t.Errorf("stopPipeline error : %v", err)
	}
	t.Log("the status code is ", code)
	t.Log("the response body :", string(body))
	t.Log("the status code is ", code1)
	t.Log("the response body :", string(body1))
}
