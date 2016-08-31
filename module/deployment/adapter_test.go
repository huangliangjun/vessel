package deployment

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/setting"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	setting.RunTime = &setting.RunTimeConf{}
	setting.RunTime.K8s.Host = "127.0.0.1"
	setting.RunTime.K8s.Port = "8080"
}

func TestDeploy(t *testing.T) {
	data, err := ioutil.ReadFile("./stage.json")
	if err != nil {
		t.Error(err)
		return
	}

	stage := &models.Stage{}
	err = json.Unmarshal(data, stage)
	if err != nil {
		t.Error(err)
		return
	}

	deployment := NewDeployment(stage)
	log.Println(deployment)
	result := deployment.Deploy()
	log.Println(result)
}
