package models

import (
	"fmt"
	"time"

	"github.com/containerops/vessel/setting"
	"github.com/coreos/etcd/client"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned"
)

var (
	// EtcdClient etcd client
	EtcdClient client.Client
	// K8sClient k8s client
	K8sClient *unversioned.Client
)

const (
	// EtcdConnectPath connect path for etcd
	EtcdConnectPath = "http://%s:%s"
	// K8sConnectPath connect path for k8s
	K8sConnectPath = "%v:%v"
)

// InitEtcd for etcd init
func InitEtcd() error {
	if EtcdClient == nil {
		var etcdEndPoints []string
		for _, value := range setting.RunTime.Etcd.Endpoints {
			etcdEndPoints = append(etcdEndPoints, fmt.Sprintf(EtcdConnectPath, value["host"], value["port"]))
		}

		cfg := client.Config{
			Endpoints: etcdEndPoints,
			Transport: client.DefaultTransport,
			// Set timeout per request to fail fast when the target endpoint is unavailable
			HeaderTimeoutPerRequest: time.Second,
		}
		var err error
		EtcdClient, err = client.New(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

// InitK8S for K8S init
func InitK8S() error {
	if K8sClient == nil {
		clientConfig := restclient.Config{}
		host := fmt.Sprintf(K8sConnectPath, setting.RunTime.K8s.Host, setting.RunTime.K8s.Port)
		clientConfig.Host = host
		// ClientConfig.Host = setting.RunTime.Database.Host
		client, err := unversioned.New(&clientConfig)
		if err != nil {
			fmt.Printf("New unversioned client err: %v!\n", err.Error())
			return err
		}
		K8sClient = client
	}
	return nil
}
