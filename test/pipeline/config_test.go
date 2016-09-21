package pipeline

import (
	"testing"

	"github.com/astaxie/beego/config"
)

var (
	RequestUrl string
	ListenType string
)

func Test_GetVesselConf(t *testing.T) {
	path := "../testunit.conf"
	conf, err := config.NewConfig("ini", path)
	if err != nil {
		t.Errorf("Read %v error : %v", path, err)
	}
	if requestUrl := conf.String("test::request_url"); requestUrl != "" {
		RequestUrl = requestUrl
	} else {
		t.Errorf("Read request_url error ")
	}
	if listenType := conf.String("test::listen_type"); listenType != "" {
		ListenType = listenType
	} else {
		t.Errorf("Read listen_type error ")
	}
}
