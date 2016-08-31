package pipeline

import (
	"encoding/json"
	"testing"

	"github.com/containerops/vessel/models"
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
                "type": "Check",
                "triggers": "frontend1",
                "conditions": "redis-slave,frontend"
            },
            {
                "type": "End",
                "triggers": "",
                "conditions": "frontend2"
            }
        ],
        "Stages": [
            {
                "Name": "redis-master",
                "Type": "container",
				"Replicas":1,
                "Dependencies": "",
                "Artifacts": [
                    {
                        "Name": "aaa",
                        "Path": "/a"
                    }
                ],
                "Volumes": [
                    {
                        "Name": "aaa",
                        "HostPath": "/a"
                    }
                ]
            },
			{
                "Name": "redis-slave",
                "Type": "container",
				"Replicas":2,
                "Dependencies": "redis-master",
                "Artifacts": [
                    {
                        "Name": "bbb",
                        "Path": "/b"
                    }
                ],
                "Volumes": [
                    {
                        "Name": "bbb",
                        "HostPath": "/b"
                    }
                ]
            },
			{
                "Name": "frontend",
                "Type": "container",
				"Replicas":1,
                "Dependencies": "redis-master",
                "Artifacts": [
                    {
                        "Name": "ccc",
                        "Path": "/c"
                    }
                ],
                "Volumes": [
                    {
                        "Name": "ccc",
                       "HostPath": "/c"
                    }
                ]
            },
			{
                "Name": "frontend1",
                "Type": "container",
				"Replicas":1,
                "Dependencies": "frontend",
                "Artifacts": [
                    {
                        "Name": "ddd",
                        "Path": "/d"
                    }
                ],
                "Volumes": [
                    {
                        "Name": "ddd",
                        "HostPath": "/d"
                    }
                ]
            },
			{
                "Name": "frontend2",
                "Type": "container",
				"Replicas":1,
                "Dependencies": "frontend1",
                "Artifacts": [
                    {
                        "Name": "eee",
                        "Path": "/e"
                    }
                ],
                "Volumes": [
                    {
                        "Name": "eee",
                        "HostPath": "/e"
                    }
                ]
            }
			
        ]
    }
}`

func Test_CreatePipeline(t *testing.T) {
	var pipelineTemplate models.PipelineTemplate
	err := json.Unmarshal([]byte(pipelineJson), &pipelineTemplate)
	if err != nil {
		t.Error("json Unmarshal failure")
	} else {
		t.Log(pipelineTemplate)
	}

	t.Log(string(CreatePipeline(&pipelineTemplate)))

}

func Test_StartPipeline(t *testing.T) {
	var pid uint64 = 1
	t.Log(string(StartPipeline(pid)))

}
