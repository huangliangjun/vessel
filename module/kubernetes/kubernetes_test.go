package kubernetes

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/utils/timer"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned"
)

func TestCreateInK8s(t *testing.T) {
	if err := createClient("127.0.0.1", "8080"); err != nil {
		t.Errorf("Create client err : %v", err.Error())
		return
	}
	stage := getStage()
	hourglass := timer.InitHourglass(time.Duration(20))
	res := CreateStage(stage, hourglass)
	if res.Result != models.ResultSuccess {
		t.Errorf("Create stage err : %v", res.Detail)
		return
	}
	res = GetBusinessRes(stage, hourglass)
	if res.Result != models.ResultSuccess {
		t.Errorf("Create stage err : %v", res.Detail)
		return
	}
}

func TestDeleteInK8s(t *testing.T) {
	if err := createClient("127.0.0.1", "8080"); err != nil {
		t.Errorf("Create client err : %v", err.Error())
		return
	}
	res := DeleteStage(getStage(), timer.InitHourglass(time.Duration(20)))
	if res.Result != models.ResultSuccess {
		t.Errorf("Delete stage err : %v", res.Detail)
	}
}

func createClient(host string, prop string) error {
	if k8sClient == nil {
		clientConfig := restclient.Config{
			Host: fmt.Sprintf(models.K8sConnectPath, host, prop),
		}
		var err error
		k8sClient, err = unversioned.New(&clientConfig)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func getStage() *models.Stage {
	return &models.Stage{
		Name:                "etcdStage",
		Namespace:           "xxxx",
		Replicas:            3,
		Image:               "unknow",
		Port:                80,
		StatusCheckURL:      "/heath",
		StatusCheckInterval: 30,
		StatusCheckCount:    3,
		EnvName:             "",
		EnvValue:            "",
		Dependence:          "stageA,stageB,stageC",
		Status:              models.StateRunning,
	}
}
