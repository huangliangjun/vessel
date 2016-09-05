package deployment

import (
	"bytes"
	"errors"
	"fmt"
	"log"

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

func DeployInK8S(data *K8SData) (detail string, err error) {
	client, err := kubernetes.NewRESTClient(fmt.Sprintf("%s:%s", setting.RunTime.K8s.Host, setting.RunTime.K8s.Port))
	if err != nil {
		detail = err.Error()
		return
	}

	if data.Namespace != "" {
		detail, err = newNamespace(client, data.Namespace)
		if err != nil {
			return detail, err
		}
	}

	detail, err = createRC(client, data)
	if err != nil {
		return detail, err
	}

	return createService(client, data)
	//TODO close client
}

//Check namespace, create one if not exist
func newNamespace(client *kubernetes.RESTClient, namespace string) (string, error) {
	params := kubernetes.NewParamsWithResourceType(kubernetes.NAMESPACES, namespace, "", false, false)
	result := client.Get(params)
	if result.Err != nil {
		return result.Err.Error(), result.Err
	}
	log.Println("Get Namespace Code: ", result.StatusCode)

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
		log.Println("Create Namespace Code: ", result.StatusCode)

		if result.StatusCode == K8S_CREATED {
			return string(result.Body), nil
		}
	}

	return string(result.Body), errors.New(fmt.Sprintf("Respond code is: ", result.StatusCode))
}

func createRC(client *kubernetes.RESTClient, data *K8SData) (string, error) {
	params := kubernetes.NewParamsWithResourceType(kubernetes.REPLICATIONCONTROLLERS, data.Name, data.Namespace, false, false)
	meta, err := params.EncodingParams()
	if err != nil {
		return err.Error(), err
	}

	buffer := make([][]byte, 4)
	buffer[2], err = data.EncodingData(kubernetes.REPLICATIONCONTROLLERS)
	if err != nil {
		return err.Error(), err
	}
	buffer[0] = bytes.TrimRight(meta, "}")
	buffer[1] = []byte("},\"spec\":")
	buffer[3] = []byte("}")
	body := bytes.Join(buffer, []byte(""))
	log.Println(string(body))

	result := client.Create(params, body)
	if result.Err != nil {
		return result.Err.Error(), result.Err
	}
	log.Println("Create ReplicationController Code: ", result.StatusCode)

	if result.StatusCode == K8S_CREATED {
		return string(result.Body), nil
	}

	return string(result.Body), errors.New(fmt.Sprintf("Respond code is: ", result.StatusCode))
}

func createService(client *kubernetes.RESTClient, data *K8SData) (string, error) {
	params := kubernetes.NewParamsWithResourceType(kubernetes.SERVICES, data.Name, data.Namespace, false, false)
	meta, err := params.EncodingParams()
	if err != nil {
		return err.Error(), err
	}

	buffer := make([][]byte, 4)
	buffer[2], err = data.EncodingData(kubernetes.SERVICES)
	if err != nil {
		return err.Error(), err
	}
	buffer[0] = bytes.TrimRight(meta, "}")
	buffer[1] = []byte("},\"spec\":")
	buffer[3] = []byte("}")
	body := bytes.Join(buffer, []byte(""))
	log.Println(string(body))

	result := client.Create(params, body)
	if result.Err != nil {
		return result.Err.Error(), result.Err
	}
	log.Println("Create Service Code: ", result.StatusCode)

	if result.StatusCode == K8S_CREATED {
		return string(result.Body), nil
	}

	return string(result.Body), errors.New(fmt.Sprintf("Respond code is: ", result.StatusCode))
}

func UndeployInK8S(data *K8SData) (detail string, err error) {
	client, err := kubernetes.NewRESTClient(fmt.Sprintf("%s:%s", setting.RunTime.K8s.Host, setting.RunTime.K8s.Port))
	if err != nil {
		detail = err.Error()
		return
	}

	detail, err = deleteRC(client, data)
	if err != nil {
		return detail, err
	}

	return deleteService(client, data)
	//TODO close client
}

func deleteRC(client *kubernetes.RESTClient, data *K8SData) (string, error) {
	params := kubernetes.NewParamsWithResourceType(kubernetes.REPLICATIONCONTROLLERS, data.Name, data.Namespace, false, false)

	patch := []byte(`[{"op":"replace","path":"/spec/replicas","value":0}]`)
	result := client.Update(params, patch)
	if result.Err != nil {
		return result.Err.Error(), result.Err
	}
	if result.StatusCode != K8S_OK {
		return string(result.Body), errors.New(fmt.Sprintf("Respond code is: ", result.StatusCode))
	}

	result = client.Delete(params)
	if result.Err != nil {
		return result.Err.Error(), result.Err
	}
	log.Println("Delete ReplicationController Code: ", result.StatusCode)

	if result.StatusCode == K8S_OK {
		return string(result.Body), nil
	}

	return string(result.Body), errors.New(fmt.Sprintf("Respond code is: ", result.StatusCode))
}

func deleteService(client *kubernetes.RESTClient, data *K8SData) (string, error) {
	params := kubernetes.NewParamsWithResourceType(kubernetes.SERVICES, data.Name, data.Namespace, false, false)

	result := client.Delete(params)
	if result.Err != nil {
		return result.Err.Error(), result.Err
	}
	log.Println("Delete Service Code: ", result.StatusCode)

	if result.StatusCode == K8S_OK {
		return string(result.Body), nil
	}

	return string(result.Body), errors.New(fmt.Sprintf("Respond code is: ", result.StatusCode))
}
