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
			            "path": "gcr.io/google_containers/gb-redisslave:v1",
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
                "Type": "container",
				"Replicas":1,
                "Dependencies": "redis-slave",
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

func Test_StopPipeline(t *testing.T) {
	var pid uint64 = 1
	var pvid uint64 = 1
	t.Log(string(StopPipeline(pid, pvid)))

}

func Test_DeletePipeline(t *testing.T) {
	var pid uint64 = 1
	t.Log(string(DeletePipeline(pid)))

}
