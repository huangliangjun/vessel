package kubernetes

import (
	//"io/ioutil"
	"log"
	"testing"

	//"k8s.io/kubernetes/pkg/api/v1"
)

const (
	HOST = "127.0.0.1:8080"
)

var client *RESTClient

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	client, _ = NewRESTClient(HOST)
}

func TestNamespace(t *testing.T) {
	params := NewParamsWithResourceType(NAMESPACES, "vessel", "", false, false)
	log.Println(params)

	body, err := params.EncodingParams()
	if err != nil {
		t.Error(err)
		return
	}

	result := client.Create(params, body)
	log.Println(result.StatusCode)
	if result.Err != nil {
		t.Error(result.Err)
		return
	}
	log.Println(string(result.Body))

	result = client.Get(params)
	log.Println(result.StatusCode)
	if result.Err != nil {
		t.Error(result.Err)
		return
	}
	log.Println(string(result.Body))

	result = client.Delete(params)
	log.Println(result.StatusCode)
	if result.Err != nil {
		t.Error(result.Err)
		return
	}
	log.Println(string(result.Body))
}

/*
func PostGetDeleteResourceWithJson(t *testing.T, fileName string) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Error(err)
		return
	}

	params, err := NewParams(data)
	if err != nil {
		t.Error(err)
		return
	}

	result := client.Create(params)
	log.Println(result)

	result = client.Get(params)
	log.Println(result)

	result = client.Delete(params)
	log.Println(result)
}

func TestResource(t *testing.T) {
	PostGetDeleteResourceWithJson(t, "./testpod.yaml")
	PostGetDeleteResourceWithJson(t, "./testrc.yaml")
}
*/
