package deployment

import (
	"log"
	"strings"

	"github.com/containerops/vessel/models"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/util/intstr"
)

type Deployment struct {
	ID         uint64
	Name       string
	Namespace  string
	DeployType string
	*DeployData
}

func NewDeployment(stage *models.Stage) *Deployment {
	data := &DeployData{}
	replicas := new(int32)
	switch stage.Type {
	case models.STAGECONTAINER:
		*replicas = int32(stage.Replicas)
		data.K8S = &K8SData{
			Name:        stage.Name,
			Namespace:   stage.Namespace,
			Replicas:    replicas,
			Labels:      map[string]string{"app": stage.Name},
			PodSpec:     newPodSpec(stage.Artifacts, stage.Volumes),
			ServiceSpec: newServiceSpec(stage.Ports, stage.Artifacts[0].Container.Ports),
		}
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

func newPodSpec(a []models.Artifact, v []models.Volume) *v1.PodSpec {
	if len(a) == 0 || len(a) != len(v) {
		log.Println("Length is: ", len(a))
		return nil
	}

	size := len(a)
	containers := make([]v1.Container, size)
	volumes := make([]v1.Volume, size)
	for i := 0; i < size; i++ {
		containers[i] = v1.Container{
			Name:  a[i].Name,
			Image: a[i].Path,
			Command: func() []string {
				if a[i].Lifecycle != nil && a[i].Lifecycle.Runtime != nil {
					return a[i].Lifecycle.Runtime
				} else {
					return nil
				}
			}(),
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
			Lifecycle: func() *v1.Lifecycle {
				if a[i].Lifecycle == nil {
					return nil
				}

				var postStart *v1.Handler
				var preStop *v1.Handler

				if a[i].Lifecycle.Before != nil {
					postStart = &v1.Handler{
						Exec: &v1.ExecAction{
							Command: a[i].Lifecycle.Before,
						},
					}
				}

				if a[i].Lifecycle.After != nil {
					preStop = &v1.Handler{
						Exec: &v1.ExecAction{
							Command: a[i].Lifecycle.After,
						},
					}
				}

				if postStart != nil || preStop != nil {
					return &v1.Lifecycle{
						PostStart: postStart,
						PreStop:   preStop,
					}
				} else {
					return nil
				}
			}(),
		}

		volumes[i] = v1.Volume{
			Name: v[i].Name,
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: v[i].HostPath,
				},
			},
		}
	}

	return &v1.PodSpec{
		Volumes:    volumes,
		Containers: containers,
	}
}

func newServiceSpec(s []models.ServicePort, c []models.ContainerPort) *v1.ServiceSpec {
	size := len(s)
	if size == 0 {
		log.Println("Length is: ", len(s))
		return nil
	}

	ports := make([]v1.ServicePort, size)
	for i := 0; i < size; i++ {
		ports[i] = v1.ServicePort{
			Name: s[i].Name,
			Port: s[i].Port,
			TargetPort: func() intstr.IntOrString {
				for _, v := range c {
					if v.Name == s[i].Name {
						return intstr.FromInt(int(v.ContainerPort))
					}
				}
				return intstr.FromInt(int(s[i].Port))
			}(),
		}
	}
	return &v1.ServiceSpec{
		Ports: ports,
	}
}

func (d *Deployment) Deploy() *models.StageResult {
	switch d.DeployType {
	case models.STAGECONTAINER:
		detail, err := DeployInK8S(d.K8S)
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

func (d *Deployment) Undeploy() *models.StageResult {
	switch d.DeployType {
	case models.STAGECONTAINER:
		detail, err := UndeployInK8S(d.K8S)
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
