package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/containerops/vessel/setting"
	"github.com/coreos/etcd/client"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned"
)

var (
	// EtcdClient etcd client
	EtcdClient client.Client
	// K8sClient k8s client
	K8sClient *unversioned.Client
	// Db mysql client
	Db *gorm.DB
	//the sql data is valid (general)
	DataValidStatus uint = 0
	//the sql data is invalid (delete)
	DataInValidStatus uint = 1
	ErrTxHasBegan          = errors.New("<Gorm.Begin> transaction already begin")
	ErrTxDone              = errors.New("<Gorm.Commit/Rollback> transaction not begin")
	ErrMultiRows           = errors.New("<Gorm.QuerySeter> return multi rows")
	ErrNoRows              = errors.New("<Gorm.QuerySeter> no row found")
	ErrArgs                = errors.New("<Gorm.Args> args error may be empty")

	VersionReady   = "ready"
	VersionRunning = "running"
	VersionDelete  = "Delete"
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

// InitDatabase for mysql init
func InitDatabase() error {
	fmt.Println("the Db is ", Db)
	var err error
	if Db == nil {
		dbArgs := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=%s&loc=%s&parseTime=%s",
			setting.RunTime.Database.Username, setting.RunTime.Database.Password,
			setting.RunTime.Database.Protocol, setting.RunTime.Database.Host, setting.RunTime.Database.Port,
			setting.RunTime.Database.Schema, setting.RunTime.Database.Param["charset"],
			setting.RunTime.Database.Param["loc"], setting.RunTime.Database.Param["parseTime"])
		Db, err = gorm.Open("mysql", dbArgs)

		if err != nil {
			panic(err)
		}
		Db.LogMode(true)
		Db.DB().SetMaxIdleConns(10)
		Db.DB().SetMaxOpenConns(100)
		Db.SingularTable(true)
		if err = Sync(); err != nil {
			panic(err)
		}
	}
	return nil
}

//Sync database structs
func Sync() error {
	fmt.Println("Sync database structs ")
	Db.AutoMigrate(&Pipeline{})
	Db.AutoMigrate(&PipelineVersion{})
	Db.AutoMigrate(&Point{})
	Db.AutoMigrate(&Stage{})
	Db.AutoMigrate(&StageVersion{})

	return nil
}
