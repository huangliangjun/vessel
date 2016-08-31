package deployment

import (
	"log"
	"strings"

	"github.com/containerops/vessel/models"
	"k8s.io/kubernetes/pkg/api/v1"
)

// Only one of its members may be specified.
type DeployData struct {
	K8S []K8SData
	VM  []VMData
	PC  []PCData
}

type K8SData struct {
	v1.Container
	v1.Volume
}

type VMData struct {
}

type PCData struct {
}

type Deployment struct {
	ID         uint64
	Name       string
	Namespace  string
	DeployType string
	DeployData
}

func NewDeployment(stage *models.Stage) *Deployment {
	data := DeployData{}
	switch stage.Type {
	case models.STAGECONTAINER:
		data.K8S = newK8SData(stage.Artifacts, stage.Volumes)
	case models.STAGEVM:
		return nil
	case models.STAGEPC:
		return nil
	default:
		return nil
	}

	return &Deployment{
		ID:         stage.ID,
		Name:       stage.Name,
		Namespace:  stage.Namespace,
		DeployType: stage.Type,
		DeployData: data,
	}
}

func newK8SData(a []models.Artifact, v []models.Volume) []K8SData {
	if len(a) != len(v) {
		log.Println("Length dismatch.")
		return nil
	}

	size := len(a)
	k8sData := make([]K8SData, size)
	for i := 0; i < size; i++ {
		k8sData[i].Container = v1.Container{
			Name:       a[i].Name,
			Image:      a[i].Path,
			Command:    a[i].Lifecycle.Runtime,
			WorkingDir: a[i].Container.WorkingDir,
			Ports: func(ports []models.ContainerPort) []v1.ContainerPort {
				psize := len(ports)
				parr := make([]v1.ContainerPort, psize)
				for j := 0; j < psize; j++ {
					parr[j] = v1.ContainerPort{
						Name:          ports[j].Name,
						HostPort:      ports[j].HostPort,
						ContainerPort: ports[j].ContainerPort,
					}
				}
				return parr
			}(a[i].Container.Ports),
			Env: func(envs []models.EnvVar) []v1.EnvVar {
				esize := len(envs)
				earr := make([]v1.EnvVar, esize)
				for j := 0; j < esize; j++ {
					earr[j] = v1.EnvVar{
						Name:  envs[j].Name,
						Value: envs[j].Value,
					}
				}
				return earr
			}(a[i].Container.Env),
			Lifecycle: &v1.Lifecycle{
				PostStart: &v1.Handler{
					Exec: &v1.ExecAction{
						Command: a[i].Lifecycle.Before,
					},
				},
				PreStop: &v1.Handler{
					Exec: &v1.ExecAction{
						Command: a[i].Lifecycle.After,
					},
				},
			},
		}

		k8sData[i].Volume = v1.Volume{
			Name: v[i].Name,
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: v[i].HostPath,
				},
			},
		}
	}

	return k8sData
}

func (d *Deployment) Deploy() *models.StageResult {
	switch d.DeployType {
	case models.STAGECONTAINER:
		detail, err := DeployInK8S(d.Namespace, d.K8S)
		result := &models.StageResult{
			ID:        d.ID,
			Name:      d.Name,
			Namespace: d.Namespace,
			Detail:    detail,
		}
		if err != nil {
			if strings.Contains(err.Error(), models.ResultTimeout) {
				result.Status = models.ResultTimeout
			} else {
				result.Status = models.ResultFailed
			}
		} else {
			result.Status = models.ResultSuccess
		}
		return result
	case models.STAGEVM:
		return nil
	case models.STAGEPC:
		return nil
	default:
		return nil
	}
}
