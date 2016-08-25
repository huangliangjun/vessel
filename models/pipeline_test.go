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

func Test_AddPipeline(t *testing.T) {
	var pipelineSpecTemplate PipelineSpecTemplate
	err := json.Unmarshal([]byte(pipelineJson), &pipelineSpecTemplate)
	if err != nil {
		t.Error("json Unmarshal failure")
	} else {
		t.Log(pipelineSpecTemplate.MetaData)
	}
	id, err := AddPipeline(&pipelineSpecTemplate)
	if id == 0 || err != nil {
		t.Error("AddPipeline failure ", err)
	}

}

func Test_QueryPipeline(t *testing.T) {

	where := map[string]interface{}{
		"id":        1,
		"namespace": "vessel",
		"name":      "vessel",
	}
	pipelineSpecTemplate, err := QueryPipeline(where)
	if err != nil {
		t.Error("QueryPipeline failure ", err)
	} else {
		t.Log(pipelineSpecTemplate.MetaData)
	}

}

func Test_DeletePipeline(t *testing.T) {

	var pid int64 = 1
	err := DeletePipeline(pid)
	if err != nil {
		t.Error("DeletePipeline failure ", err)
	} else {
		t.Log("DeletePipeline success")
	}

}

func Test_UpdatePipeline(t *testing.T) {
	pipelineSpecTemplate := new(PipelineSpecTemplate)
	err := UpdatePipeline(pipelineSpecTemplate)
	if err != nil {
		t.Error("UpdatePipeline failure ", err)
	} else {
		t.Log("UpdatePipeline success")
	}
}

func Test_AddPipelineVersion(t *testing.T) {
	pipelineVersion := &PipelineVersion{
		Pid:           1,
		Detail:        "",
		VersionStatus: "ready",
	}
	id, err := AddPipelineVersion(pipelineVersion)
	if id == -1 || id == 0 || err != nil {
		t.Error("AddPipelineVersion failure ", err)
	} else {
		t.Log("AddPipelineVersion success ", pipelineVersion)
	}
}

func Test_UpdatePipelineVersion(t *testing.T) {
	pipelineVersion := &PipelineVersion{
		Pid:           1,
		Detail:        "",
		VersionStatus: "running",
	}
	err := UpdatePipelineVersion(pipelineVersion)
	if err != nil {
		t.Error("UpdatePipelineVersion failure ", err)
	} else {
		t.Log("UpdatePipelineVersion success")
	}
}

func Test_QueryPipelineVersionByPid(t *testing.T) {
	var pid int64 = 1
	pipelineVersion, err := QueryPipelineVersionByPid(pid)
	if err != nil {
		t.Error("GetPipelineVersionByPid failure ", err)
	} else {
		t.Log("GetPipelineVersionByPid success ", pipelineVersion)
	}
}
