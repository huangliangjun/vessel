package deployment

import (
	"encoding/json"
	"errors"

	"github.com/containerops/vessel/module/kubernetes"
	"k8s.io/kubernetes/pkg/api/v1"
)

// Only one of its members may be specified.
type DeployData struct {
	K8S *K8SData
	VM  *VMData
	PC  *PCData
}

type K8SData struct {
	Name      string
	Namespace string
	Replicas  *int32
	Labels    map[string]string
	*v1.PodSpec
	*v1.ServiceSpec
}

func (k *K8SData) EncodingData(r kubernetes.ResourceType) ([]byte, error) {
	switch r {
	case kubernetes.REPLICATIONCONTROLLERS:
		return k.encodingRC()
	case kubernetes.SERVICES:
		return k.encodingService()
	default:
		return nil, errors.New("Unsupported Resource")
	}
}

func (k *K8SData) encodingRC() ([]byte, error) {
	if k.PodSpec == nil {
		return nil, errors.New("Empty PodSpec")
	}

	rcspec := v1.ReplicationControllerSpec{
		Replicas: k.Replicas,
		Selector: k.Labels,
		Template: &v1.PodTemplateSpec{
			ObjectMeta: v1.ObjectMeta{
				Name:      k.Name,
				Namespace: k.Namespace,
				Labels:    k.Labels,
			},
			Spec: *k.PodSpec,
		},
	}
	return json.Marshal(&rcspec)
}

func (k *K8SData) encodingService() ([]byte, error) {
	if k.ServiceSpec == nil {
		return nil, errors.New("Empty ServiceSpec")
	}

	k.ServiceSpec.Selector = k.Labels
	return json.Marshal(k.ServiceSpec)
}

type VMData struct {
}

type PCData struct {
}
