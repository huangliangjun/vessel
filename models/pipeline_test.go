package models

import (
	"encoding/json"
	"testing"
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
                "Type": "star",
                "Triggers": "a,b,c",
                "Conditions": "d,e,f"
            }
        ],
        "Stages": [
            {
                "Name": "A",
                "Type": "container",
                "Dependencies": "a,b,c",
                "Artifacts": [
                    {
                        "Name": "aaa",
                        "Path": "/d"
                    }
                ],
                "Volumes": [
                    {
                        "Name": "bbb",
                        "HostPath": {
                            "Path": "/e"
                        },
                        "EmptyDir": {
                            "Medium": "awf"
                        }
                    }
                ]
            }
        ]
    }
}`

func Test_PipelineAdd(t *testing.T) {
	var pipelineSpecTemplate PipelineSpecTemplate
	err := json.Unmarshal([]byte(pipelineJson), &pipelineSpecTemplate)
	if err != nil {
		t.Error("json Unmarshal failure")
	} else {
		t.Log(pipelineSpecTemplate.MetaData)
	}
	pipeline := pipelineSpecTemplate.MetaData

	if err = pipeline.Add(); err != nil {
		t.Error("Pipeline Add failure ", err)
	} else {
		t.Log("Pipeline Add success .")
	}

}

func Test_PipelineQueryOne(t *testing.T) {

	pipeline := &Pipeline{
		Id:        1,
		Namespace: "vessel",
		Name:      "vessel",
		//Status:    1,
	}

	if err := pipeline.QueryOne(); err != nil {
		t.Error("Pipeline QueryOne failure ", err)
	} else {
		t.Log(pipeline)
	}

}

func Test_PipelineDelete(t *testing.T) {
	pipeline := &Pipeline{
		Id: 1,
	}

	if err := pipeline.Delete(); err != nil {
		t.Error("Pipeline Delete failure ", err)
	} else {
		t.Log("Pipeline Delete success")
	}

}

func Test_PipelineUpdate(t *testing.T) {
	pipeline := &Pipeline{}

	if err := pipeline.Update(); err != nil {
		t.Error("Update Pipeline failure ", err)
	} else {
		t.Log("Update Pipeline success")
	}
}

func Test_AddPipelineVersion(t *testing.T) {
	pv := &PipelineVersion{
		Pid:           1,
		Detail:        "",
		VersionStatus: "ready",
	}

	if err := pv.Add(); err != nil {
		t.Error("PipelineVersion Add failure ", err)
	} else {
		t.Log("PipelineVersion Add success ", pv)
	}
}

func Test_PipelineVersionUpdate(t *testing.T) {
	pv := &PipelineVersion{
		Pid:           1,
		Detail:        "",
		VersionStatus: "running",
	}

	if err := pv.Update(); err != nil {
		t.Error("PipelineVersion Update failure ", err)
	} else {
		t.Log("PipelineVersion Update success", pv)
	}
}

func Test_PipelineVersionQueryOne(t *testing.T) {
	pv := &PipelineVersion{
		Pid: 1,
	}

	if err := pv.QueryOne(); err != nil {
		t.Error("PipelineVersion QueryOne failure ", err)
	} else {
		t.Log("PipelineVersion QueryOne success ", pv)
	}
}
