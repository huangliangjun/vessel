package pipeline

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

var (
	Get     string = "GET"
	Post    string = "POST"
	Put     string = "PUT"
	patch   string = "PATCH"
	Delete  string = "DELETE"
	Head    string = "HEAD"
	Options string = "OPTIONS"
)

func init() {

}

func postPipeline(body io.Reader) ([]byte, string, error) {
	url := ListenType + "://" + RequestUrl + "/vessel/"
	header := map[string]string{
		"Content-Type": "application/json",
	}
	resp, err := sendHttpRequest(Post, url, body, header)
	defer resp.Body.Close()
	message, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return message, resp.Status, nil
}

func deletePipeline(pid uint64) ([]byte, string, error) {
	pidStr := strconv.FormatUint(pid, 10)
	url := ListenType + "://" + RequestUrl + "/vessel/" + pidStr + "/"
	header := map[string]string{}
	resp, err := sendHttpRequest(Delete, url, nil, header)
	defer resp.Body.Close()
	message, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return message, resp.Status, nil
}

func startPipeline(pid uint64) ([]byte, string, error) {
	pidStr := strconv.FormatUint(pid, 10)
	url := ListenType + "://" + RequestUrl + "/vessel/" + pidStr + "/"
	header := map[string]string{}
	resp, err := sendHttpRequest(Post, url, nil, header)
	defer resp.Body.Close()
	message, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return message, resp.Status, nil
}

func stopPipeline(pid, pvid uint64) ([]byte, string, error) {
	pidStr := strconv.FormatUint(pid, 10)
	pvidStr := strconv.FormatUint(pvid, 10)
	url := ListenType + "://" + RequestUrl + "/vessel/" + pidStr + "/" + pvidStr + "/"
	header := map[string]string{}
	resp, err := sendHttpRequest(Delete, url, nil, header)
	defer resp.Body.Close()
	message, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return message, resp.Status, nil
}

func sendHttpRequest(method, requestUrl string, body io.Reader, header map[string]string) (*http.Response, error) {
	url, err := url.Parse(requestUrl)
	if err != nil {
		return &http.Response{}, err
	}
	var client *http.Client
	switch url.Scheme {
	case "":
		fallthrough
	case "https":
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	case "http":
		client = &http.Client{}
	default:
		return &http.Response{}, fmt.Errorf("bad url schema: %v", url.Scheme)
	}
	req, err := http.NewRequest(method, url.String(), body)

	for k, v := range header {
		req.Header.Set(k, v)
	}
	if err != nil {
		return &http.Response{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("the error :", err)
		return &http.Response{}, err
	}
	return resp, nil
}
