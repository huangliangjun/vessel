package stage

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/setting"
	"github.com/containerops/vessel/utils/timer"
)

func TestStartStage(t *testing.T) {
	if err := setting.InitConf("./conf/global.yaml", "./conf/runtime.yaml"); err != nil {
		t.Logf("Read config error: %v", err.Error())
		return
	}
	if err := models.InitEtcd(); err != nil {
		t.Logf("Init etcd error: %v", err.Error())
		return
	}
	if err := models.InitK8S(); err != nil {
		t.Logf("Init k8s error: %v", err.Error())
		return
	}
	str := jsonStr()
	pipelineTemp := &models.PipelineSpecTemplate{}
	err := json.Unmarshal([]byte(str), pipelineTemp)
	if err != nil {
		t.Log(err)
		return
	}

	count := len(pipelineTemp.Spec)
	finishChan := make(chan *models.ExecutedResult, count)
	resultChan := make(chan []interface{})
	hourglass := timer.InitHourglass(int64(60) * int64(time.Second))
	go isFinish(finishChan, resultChan, count)
	for _, stageSpec := range pipelineTemp.Spec {
		stageSpec.Namespace = pipelineTemp.MetaData.Namespace
		stageSpec.PipelineName = pipelineTemp.MetaData.Name
		go StartStage(stageSpec, finishChan, hourglass)
	}
	t.Log(<-resultChan)
}

func isFinish(finishChan chan *models.ExecutedResult, resultChan chan []interface{}, count int64) {
	resultList := make([]interface{}, 0, count)
	running := true
	for running {
		select {
		case result := <-finishChan:
			resultList = append(resultList, result)
			if result.Status != models.ResultSuccess {
				running = false
			} else {
				running = len(resultList) != count
			}
		}
	}
	resultChan <- resultList
}

func jsonStr() string {
	return `{
	    "kind": "CCloud",
	    "apiVersion": "v1",
	    "status": "",
	    "apiServerUrl": "",
	    "apiServerAuth": "",
	    "metadata": {
		"name": "guestbook",
		"namespace": "guestbook",
		"selfLink": "",
		"uid": "",
		"creationTimestamp": "",
		"deletionTimestamp": "",
		"labels": {
		    "app": "zenlin"
		},
		"annotations": {},
		"timeoutDuration": 60
	    },
	    "spec": [{
		"name": "redis-master",
		"replicas": 1,
		"dependence": "",
		"kind": "value",
		"statusCheckLink": "/health",
		"statusCheckInterval": 0,
		"statusCheckCount": 0,
		"image": "gcr.io/google_containers/redis:e2e",
		"port": 6379,
		"envName": "",
		"envValue": ""
	    }, {
		"name": "redis-slave",
		"replicas": 2,
		"dependence": "redis-master",
		"kind": "value",
		"statusCheckLink": "/health",
		"statusCheckInterval": 0,
		"statusCheckCount": 0,
		"image": "gcr.io/google_samples/gb-redisslave:v1",
		"port": 6379,
		"envName": "",
		"envValue": ""
	    }, {
		"name": "frontend",
		"replicas": 3,
		"dependence": "redis-slave",
		"kind": "value",
		"statusCheckLink": "/health",
		"statusCheckInterval": 0,
		"statusCheckCount": 0,
		"image": "gcr.io/google_samples/gb-frontend:v3",
		"port": 80,
		"envName": "GET_HOSTS_FROM",
		"envValue": "dns"
	    }]
	}`
}
