package deployment

import (
	"errors"
	"fmt"

	"github.com/containerops/vessel/module/kubernetes"
	"github.com/containerops/vessel/setting"
)

const (
	K8S_OK                  = 200
	K8S_CREATED             = 201
	K8S_NOCONTENT           = 204
	K8S_REDIRECT            = 307
	K8S_BADREQUEST          = 400
	K8S_UNAUTHORIZED        = 401
	K8S_FORBIDDEN           = 403
	K8S_NOTFOUND            = 404
	K8S_METHODNOTALLOWED    = 405
	K8S_CONFLICT            = 409
	K8S_UNPROCESSABLEENTITY = 422
	K8S_TOOMANYREQUESTS     = 429
	K8S_INTERNALSERVERERROR = 500
	K8S_SERVICEUNAVAILABLE  = 503
	K8S_SERVERTIMEOUT       = 504
)

func DeployInK8S(namespace string, data []K8SData) (detail string, err error) {
	client, err := kubernetes.NewRESTClient(fmt.Sprintf("%s:%s", setting.RunTime.K8s.Host, setting.RunTime.K8s.Port))
	if err != nil {
		detail = err.Error()
		return
	}

	if namespace != "" {
		detail, err = newNamespace(client, namespace)
		if err != nil {
			return detail, err
		}
	} else {
		detail = ""
	}

	return
}

//Check namespace, create one if not exist
func newNamespace(client *kubernetes.RESTClient, namespace string) (string, error) {
	params := kubernetes.NewParamsWithResourceType(kubernetes.NAMESPACES, namespace, "", false, false)
	result := client.Get(params)
	if result.Err != nil {
		return result.Err.Error(), result.Err
	}

	if result.StatusCode == K8S_OK {
		return string(result.Body), nil
	}

	if result.StatusCode == K8S_NOTFOUND {
		body, err := params.EncodingParams()
		if err != nil {
			return err.Error(), err
		}

		result = client.Create(params, body)
		if result.Err != nil {
			return result.Err.Error(), result.Err
		}
		if result.StatusCode == K8S_CREATED {
			return string(result.Body), nil
		}
	}

	return string(result.Body), errors.New(fmt.Sprintf("Respond code is: ", result.StatusCode))
}
