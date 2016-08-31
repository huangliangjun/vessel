package kubernetes

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/watch"
)

const (
	_APIPATH = "/api/v1"
)

const ()

type RESTClient struct {
	*restclient.RESTClient
}

type Result struct {
	Body       []byte
	Err        error
	StatusCode int
}

func NewRESTClient(host string) (*RESTClient, error) {
	log.Println(host)
	baseURL := &url.URL{
		Scheme: "http",
		Host:   host,
	}

	config := &restclient.Config{
		ContentConfig: restclient.ContentConfig{
			NegotiatedSerializer: api.Codecs,
		},
	}
	if err := restclient.SetKubernetesDefaults(config); err != nil {
		log.Println(err)
		return nil, err
	}

	client, err := restclient.NewRESTClient(baseURL, _APIPATH, config.ContentConfig, config.QPS, config.Burst, config.RateLimiter, &http.Client{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &RESTClient{client}, nil
}

func (c *RESTClient) Get(params *Params) *Result {
	path := params.BuildPath()
	if path == "" {
		return &Result{
			Err:        errors.New("Not supported yet"),
			StatusCode: -1,
		}
	}

	url := c.Verb("Get").URL()
	url.Path = _APIPATH + path
	log.Println(url.String())
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return &Result{
			Err:        err,
			StatusCode: -1,
		}
	}

	return c.doRequest(req)
}

func (c *RESTClient) Create(params *Params, body []byte) *Result {
	path := params.BuildPathForPost()
	if path == "" {
		return &Result{
			Err:        errors.New("Not supported yet"),
			StatusCode: -1,
		}
	}

	url := c.Verb("Get").URL()
	url.Path = _APIPATH + path
	log.Println(url.String())
	req, err := http.NewRequest("POST", url.String(), bytes.NewReader(body))
	if err != nil {
		return &Result{
			Err:        err,
			StatusCode: -1,
		}
	}

	return c.doRequest(req)
}

func (c *RESTClient) Delete(params *Params) *Result {
	path := params.BuildPath()
	if path == "" {
		return &Result{
			Err:        errors.New("Not supported yet"),
			StatusCode: -1,
		}
	}

	url := c.Verb("Get").URL()
	url.Path = _APIPATH + path
	log.Println(url.String())
	req, err := http.NewRequest("DELETE", url.String(), nil)
	if err != nil {
		return &Result{
			Err:        err,
			StatusCode: -1,
		}
	}

	return c.doRequest(req)
}

func (c *RESTClient) doRequest(req *http.Request) *Result {
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Client.Do(req)
	if err != nil {
		return &Result{
			Err:        err,
			StatusCode: -1,
		}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Result{
			Err:        err,
			StatusCode: resp.StatusCode,
		}
	}
	return &Result{
		Body:       body,
		Err:        nil,
		StatusCode: resp.StatusCode,
	}
}

func (c *RESTClient) Watch(params *Params, events chan *watch.Event) {
	if !params.IsSetWatcher {
		params.IsSetWatcher = true
	}

	req := c.RESTClient.Get()
	log.Println(req.URL().Path)

}
