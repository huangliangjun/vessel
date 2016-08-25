package main

import (
	//"encoding/json"
	"testing"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/setting"
)

func Test_ParseConf(t *testing.T) {

	if err := setting.InitGlobalConf("./conf/global.yaml"); err != nil {
		t.Log("the parseConf err is ", err)
	}
	t.Log(setting.RunTime)
	if err := models.InitDatabase(); err != nil {
		t.Log("the parseConf err is ", err)
	}

}
