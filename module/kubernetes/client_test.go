package kubernetes

import (
	"io/ioutil"
	"log"
	"testing"
)

const (
	HOST = "127.0.0.1:8080"
)

var client *RESTClient

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	client, _ = NewRESTClient(HOST)
}

func PostResourceWithJson(fileName string) (params *Params, err error) {
	body, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}

	params, err = NewParamsWithJson(body)
	if err != nil {
		return
	}

	result := client.Create(params, body)
	log.Println("Create replicationcontrollers: ", result.StatusCode)
	return

	/*result = client.Get(params)
	log.Println(result)

	result = client.Delete(params)
	log.Println(result)*/
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

func TestPatch(t *testing.T) {
	params, err := PostResourceWithJson("./testrc.json")
	if err != nil {
		t.Error(err)
		return
	}

	patch := []byte(`[{"op":"replace","path":"/spec/replicas","value":0}]`)
	result := client.Update(params, patch)
	log.Println(string(result.Body))
	log.Println("Update replicationcontrollers: ", result.StatusCode)
	return
}

/*
func TestWatch(t *testing.T) {
	params, err := PostResourceWithJson("./testrc.json")
	if err != nil {
		t.Error(err)
		return
	}

	watcher, err := client.Watch(params)
	if err != nil {
		t.Error(err)
		return
	}

	for {
		select {
		case event := <-watcher.ResultChan():
			log.Println(event.Type)
			log.Println(event.Object.GetObjectKind())
		}
	}
}
*/
