package pipeline

import (
	"encoding/json"
	"testing"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/setting"
)

func TestStartPipeline(t *testing.T) {
	if err := setting.InitGlobalConf("./conf/global.yaml"); err != nil {
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
	bytes := StartPipeline(pipelineTemp)
	t.Log(string(bytes))
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
