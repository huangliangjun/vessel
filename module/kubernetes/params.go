package kubernetes

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type ResourceType int

const (
	BINDINGS = iota
	COMPONENTSTATUSES
	ENDPOINTS
	EVENTS
	LIMITRANGES
	NAMESPACES
	NODES
	PERSISTENTVOLUMECLAIMS
	PERSISTENTVOLUMES
	PODS
	PODTEMPLATES
	REPLICATIONCONTROLLERS
	RESOURCEQUOTAS
	SECRETS
	SERVICEACCOUNTS
	SERVICES
	UNKNOWN
)

//TODO more resource type needs to add
func NewResourceType(kind string) ResourceType {
	switch kind {
	case "Binding":
		return BINDINGS
	case "Namespace":
		return NAMESPACES
	case "Node", "NodeLists":
		return NODES
	case "Pod", "PodLists":
		return PODS
	case "ReplicationController":
		return REPLICATIONCONTROLLERS
	case "Service", "ServiceLists":
		return SERVICES
	default:
		return UNKNOWN
	}
}

func (r ResourceType) GetType() string {
	switch r {
	case BINDINGS:
		return "bindings"
	case COMPONENTSTATUSES:
		return "componentstatuses"
	case ENDPOINTS:
		return "endpoints"
	case EVENTS:
		return "events"
	case LIMITRANGES:
		return "limitranges"
	case NAMESPACES:
		return "namespaces"
	case NODES:
		return "nodes"
	case PERSISTENTVOLUMECLAIMS:
		return "persistentvolumeclaims"
	case PERSISTENTVOLUMES:
		return "persistentvolumes"
	case PODS:
		return "pods"
	case PODTEMPLATES:
		return "podtemplates"
	case REPLICATIONCONTROLLERS:
		return "replicationcontrollers"
	case RESOURCEQUOTAS:
		return "resourcequotas"
	case SECRETS:
		return "secrets"
	case SERVICEACCOUNTS:
		return "serviceaccounts"
	case SERVICES:
		return "services"
	default:
		return "unknown"
	}
}

type Params struct {
	ResourceType
	Name         string
	Namespace    string
	IsVisitProxy bool
	IsSetWatcher bool
}

type RawParams struct {
	Kind       string          `json:"kind"`
	APIVersion string          `json:"apiVersion"`
	Metadata   *ParamsMetaData `json:"metadata,omitempty"`
}

type ParamsMetaData struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

func NewParamsWithResourceType(resource ResourceType, name, namespace string, isProxy, isWatcher bool) *Params {
	return &Params{
		ResourceType: resource,
		Name:         name,
		Namespace:    namespace,
		IsVisitProxy: isProxy,
		IsSetWatcher: isWatcher,
	}
}

func NewParamsWithJson(jstr []byte) (*Params, error) {
	var r RawParams
	err := json.Unmarshal(jstr, &r)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if r.APIVersion != "v1" {
		err = fmt.Errorf("Not supported version: ", r.APIVersion)
		log.Println(err)
		return nil, err
	}

	return &Params{
		ResourceType: NewResourceType(r.Kind),
		Name:         r.Metadata.Name,
		Namespace:    r.Metadata.Namespace,
		IsVisitProxy: false,
		IsSetWatcher: false,
	}, nil
}

func (p *Params) EncodingParams() ([]byte, error) {
	raw := &RawParams{
		Kind: func(s string) string {
			switch s {
			case "replicationcontrollers":
				return "ReplicationController"
			default:
				s = strings.TrimRight(s, "s")
				upper := strings.ToUpper(s)
				return upper[:1] + s[1:]
			}
		}(p.GetType()),
		APIVersion: "v1",
		Metadata: func() *ParamsMetaData {
			if p.Name == "" && p.Namespace == "" {
				return nil
			} else {
				return &ParamsMetaData{
					Name:      p.Name,
					Namespace: p.Namespace,
				}
			}
		}(),
	}
	return json.Marshal(raw)
}

//TODO add sub path etc.
func (p *Params) BuildPath() (path string) {
	path = p.BuildPathForPost()
	if path == "" {
		return
	}

	if p.Name != "" {
		path += fmt.Sprintf("/%s", p.Name)
	}

	log.Println(path)
	return
}

func (p *Params) BuildPathForPost() (path string) {
	if p.ResourceType == UNKNOWN {
		return
	}

	if p.IsVisitProxy {
		path = "/proxy"
	} else if p.IsSetWatcher {
		path = "/watch"
	}

	if p.Namespace != "" && p.GetType() != "namespaces" {
		path += fmt.Sprintf("/namespaces/%s", p.Namespace)
	}

	path += fmt.Sprintf("/%s", p.GetType())

	log.Println(path)
	return
}
