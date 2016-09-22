package models

import (
	//	"encoding/json"
	"testing"

	"github.com/containerops/vessel/db"
	//"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/setting"
)

var pipelineJson = `{
    "Kind": "ccloud",
    "APIVersion": "v1",
    "MetaData": {
        "Namespace": "vessel",
        "Name": "vessel",
        "Timeout": 60,
        "Points": [
            {
                "type": "Start",
                "triggers": "redis-master",
                "conditions": ""
            },
            {
                "type": "Check",
                "triggers": "redis-slave,frontend",
                "conditions": "redis-master"
            },
           
            {
                "type": "End",
                "triggers": "",
                "conditions": "frontend"
            }
        ],
        "Stages": [
            {
                "Name": "redis-master",
				"Namespace":"vessel",
                "Type": "container",
				"Replicas":1,
                "Dependencies": "",
				"Ports": [
                    {
                        "name": "redis",
                        "port": 6379
                    }
                ],
                "Artifacts": [
			        {
			            "name": "redis-master",
			            "path": "gcr.io/google_containers/redis:e2e",
			            "container": {
			                "workingDir": "",
			                "ports": [
			                    {
			                        "name": "redis",
			                        "hostPort": 30001,
			                        "containerPort": 6379
			                    }
			                ],
			                "env": [
			                    {
			                        "name": "dns",
			                        "value": "redis"
			                    }
			                ]
			            }
			        }
			    ],
                "Volumes": [
                    {
                        "Name": "localvalume",
                        "HostPath": "/home/vessel"
                    }
                ]
            },
			{
                "Name": "redis-slave",
				"Namespace":"vessel",
                "Type": "container",
				"Replicas":1,
                "Dependencies": "redis-master",
				"Ports": [
                    {
                        "name": "redis",
                        "port": 6379
                    }
                ],
                "Artifacts": [
			        {
			            "name": "redis-slave",
			            "path": "gcr.io/google_samples/gb-redisslave:v1",
			            "container": {
			                "workingDir": "",
			                "ports": [
			                    {
			                        "name": "redis",
			                        "hostPort": 30002,
			                        "containerPort": 6379
			                    }
			                ],
			                "env": [
			                    {
			                        "name": "dns",
			                        "value": "redis"
			                    }
			                ]
			            }
			        }
			    ],
                "Volumes": [
                    {
                        "Name": "localvalume",
                        "HostPath": "/home/vessel"
                    }
                ]
            },
			{
                "Name": "frontend",
				"Namespace":"vessel",
                "Type": "container",
				"Replicas":1,
                "Dependencies": "redis-master",
                "Ports": [
                    {
                        "name": "frontend",
                        "port": 80
                    }
                ],
                "Artifacts": [
			        {
			            "name": "frontend",
			            "path": "gcr.io/google_samples/gb-frontend:v3",
			            "container": {
			                "workingDir": "",
			                "ports": [
			                    {
			                        "name": "frontend",
			                        "hostPort": 30003,
			                        "containerPort": 80
			                    }
			                ],
			                "env": [
			                    {
			                        "name": "dns",
			                        "value": "redis"
			                    }
			                ]
			            }
			        }
			    ],
                "Volumes": [
                    {
                        "Name": "localvalume",
                        "HostPath": "/home/vessel"
                    }
                ]
            }	
        ]
    }
}`

func Test_InitDB(t *testing.T) {
	if err := setting.InitGlobalConf("../conf/global.yaml"); err != nil {
		t.Fatal("Read config failure : ", err.Error())
	}
	if err := db.InitDB(setting.RunTime.Database.Driver, setting.RunTime.Database.Username,
		setting.RunTime.Database.Password, setting.RunTime.Database.Host+":"+setting.RunTime.Database.Port,
		setting.RunTime.Database.Schema); err != nil {
		t.Fatal(err)
	}
	if err := db.Instance.RegisterModel(new(Pipeline), new(PipelineVersion)); err != nil {
		t.Fatal(err)
	}
	if err := db.Instance.RegisterModel(new(Stage), new(StageVersion)); err != nil {
		t.Fatal(err)
	}
	if err := db.Instance.RegisterModel(new(Point), new(PointVersion)); err != nil {
		t.Fatal(err)
	}
	if err := new(Stage).AddForeignKey(); err != nil {
		fmt.Println(err)
	}
	if err := new(Point).AddForeignKey(); err != nil {
		fmt.Println(err)
	}
	if err := new(StageVersion).AddForeignKey(); err != nil {
		fmt.Println(err)
	}
	if err := new(PointVersion).AddForeignKey(); err != nil {
		fmt.Println(err)
	}
	if err := new(StageVersion).AddUniqueIndex(); err != nil {
		fmt.Println(err)
	}
	if err := new(PointVersion).AddUniqueIndex(); err != nil {
		fmt.Println(err)
	}
}

func Test_PipelineCreate(t *testing.T) {
	var pipelineTemplate PipelineTemplate
	err := json.Unmarshal([]byte(pipelineJson), &pipelineTemplate)
	if err != nil {
		t.Error("json Unmarshal failure :", err)
	} else {
		t.Log(pipelineTemplate.MetaData)
	}
	pipeline := pipelineTemplate.MetaData

	if err = pipeline.Create(); err != nil {
		t.Error("Pipeline Add failure ", err)
	} else {
		t.Log("Pipeline Add success .")
	}

}

func Test_PipelineQueryOne(t *testing.T) {

	pipeline := &Pipeline{
		//ID:        1,
		Namespace: "vessel4",
		Name:      "vessel4",
		//Status:    1,
	}

	if err := pipeline.QueryOne(); err != nil {
		t.Error("Pipeline QueryOne failure ", err)
	} else {
		t.Log(pipeline)
	}

}

func Test_PipelineUpdate(t *testing.T) {
	pipeline := &Pipeline{
		Namespace: "vessel",
		Name:      "vessel",
		Timeout:   1000,
	}

	if err := pipeline.Update(); err != nil {
		t.Error("Update Pipeline failure ", err)
	} else {
		t.Log("Update Pipeline success")
	}
}
func Test_PipelineDelete(t *testing.T) {
	pipeline := &Pipeline{
		Namespace: "vessel",
		Name:      "vessel",
	}

	if err := pipeline.SoftDelete(); err != nil {
		t.Error("Pipeline Delete failure ", err)
	} else {
		t.Log("Pipeline Delete success")
	}

}

func Test_CreatePipelineVersion(t *testing.T) {
	pv := &PipelineVersion{
		PID:    1,
		Detail: "",
		State:  "ready",
	}

	if err := pv.Create(); err != nil {
		t.Error("PipelineVersion Create failure ", err)
	} else {
		t.Log("PipelineVersion Create success ", pv)
	}
}

func Test_PipelineVersionUpdate(t *testing.T) {
	pv := &PipelineVersion{
		ID:     1,
		PID:    1,
		Detail: "",
		State:  "running",
	}

	if err := pv.Update(); err != nil {
		t.Error("PipelineVersion Update failure ", err)
	} else {
		t.Log("PipelineVersion Update success", pv)
	}
}

func Test_PipelineVersionQueryM(t *testing.T) {
	pv := &PipelineVersion{
		PID: 1,
	}
	if pvs, err := pv.QueryM(); err != nil {
		t.Error("PipelineVersion QueryOne failure ", err)
	} else {
		t.Log("PipelineVersion QueryOne success ", pvs)
	}
}

func Test_PipelineVersionSoftDelete(t *testing.T) {
	pv := &PipelineVersion{
		PID: 1,
	}
	if err := pv.SoftDelete(); err != nil {
		t.Error("PipelineVersion Delete failure ", err)
	} else {
		t.Log("PipelineVersion Delete success ")
	}
}
func Test_PipelineVersionIsExist(t *testing.T) {
	pv := &Pipeline{
		Namespace: "vessel4",
		Name:      "vessel4",
	}
	if is, err := pv.CheckIsExist(); err != nil {
		t.Error("PipelineVersionIsExist failure ", err)
	} else if is == true && err == nil {
		t.Log("PipelineVersionIsExist Success", pv)
	} else {
		t.Log("PipelineVersionIsExist record not exist ")
	}
}
