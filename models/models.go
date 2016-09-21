package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/containerops/vessel/setting"
	"github.com/coreos/etcd/client"
	_ "github.com/go-sql-driver/mysql"
	//	"github.com/jinzhu/gorm"
	"github.com/containerops/vessel/db"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned"
)

var (
	// EtcdClient etcd client
	ETCD client.Client
	// K8sClient k8s client
	K8S *unversioned.Client

	//db *gorm.DB

	//DataValidStatus the sql data is valid (general)
	DataValidStatus uint = 0
	//DataInValidStatus the sql data is invalid (delete)
	DataInValidStatus uint = 1

	ErrNotExist = errors.New("record not found")
	ErrHasExist = errors.New("record had exist")
)

const (
	// EtcdConnectPath connect path for etcd
	EtcdConnectPath = "http://%s:%s"
	// K8sConnectPath connect path for k8s
	K8sConnectPath = "%v:%v"
	// DBConnectPath connect path for DB
	DBConnectPath = "%s:%s@%s(%s:%s)/%s?charset=%s&loc=%s&parseTime=%s"
)

func init() {
	//	fmt.Println("the db is ", db)
	//	var err error
	//	if db == nil {
	//		dbArgs := "root@tcp(127.0.0.1:3306)/vesseldb?loc=Local&parseTime=True&charset=utf8"
	//		db, err = gorm.Open("mysql", dbArgs)

	//		if err != nil {
	//			panic(err)
	//		}
	//		//db.LogMode(true)
	//		db.DB().SetMaxIdleConns(10)
	//		db.DB().SetMaxOpenConns(100)
	//		db.SingularTable(true)
	//		if err = Sync(); err != nil {
	//			panic(err)
	//		}
	//	}
}

// InitEtcd for etcd init
func InitEtcd() error {
	if ETCD == nil {
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
		ETCD, err = client.New(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

// InitK8S for K8S init
func InitK8S() error {
	if K8S == nil {
		clientConfig := restclient.Config{}
		host := fmt.Sprintf(K8sConnectPath, setting.RunTime.K8s.Host, setting.RunTime.K8s.Port)
		clientConfig.Host = host
		// ClientConfig.Host = setting.RunTime.Database.Host
		client, err := unversioned.New(&clientConfig)
		if err != nil {
			log.Printf("New unversioned client err: %v!\n", err.Error())
			return err
		}
		K8S = client
	}
	return nil
}

// InitDatabase for mysql init
func InitDatabase() error {
	if err := db.InitDB(setting.RunTime.Database.Driver, setting.RunTime.Database.Username,
		setting.RunTime.Database.Password, setting.RunTime.Database.Host+":"+setting.RunTime.Database.Port,
		setting.RunTime.Database.Schema); err != nil {
		fmt.Println(err)
	}
	if err := db.Instance.RegisterModel(new(Pipeline), new(PipelineVersion)); err != nil {
		fmt.Println(err)
	}
	if err := db.Instance.RegisterModel(new(Stage), new(StageVersion)); err != nil {
		fmt.Println(err)
	}
	if err := db.Instance.RegisterModel(new(Point), new(PointVersion)); err != nil {
		fmt.Println(err)
	}
	if err := new(Stage).AddForeignKey(); err != nil {
		fmt.Println(err)
	}
	if err := new(Point).AddForeignKey(); err != nil {
		fmt.Println(err)
	}
	if err := new(StageVersion).AddForeignKey(); err != nil {
		fmt.Println(err)
	}
	if err := new(PointVersion).AddForeignKey(); err != nil {
		fmt.Println(err)
	}
	if err := new(PointVersion).AddUniqueIndex(); err != nil {
		fmt.Println(err)
	}
	return nil
}

//// InitDatabase for mysql init
//func InitDatabase() error {
//	var err error
//	if db == nil {
//		dbArgs := fmt.Sprintf(DBConnectPath, setting.RunTime.Database.Username, setting.RunTime.Database.Password,
//			setting.RunTime.Database.Protocol, setting.RunTime.Database.Host, setting.RunTime.Database.Port,
//			setting.RunTime.Database.Schema, setting.RunTime.Database.Param["charset"],
//			setting.RunTime.Database.Param["loc"], setting.RunTime.Database.Param["parseTime"])
//		if db, err = gorm.Open("mysql", dbArgs); err != nil {
//			return err
//		}
//		db.LogMode(setting.RunTime.Database.LogMode)
//		db.SingularTable(setting.RunTime.Database.SingularTable)
//		db.DB().SetMaxIdleConns(setting.RunTime.Database.MaxIdleConns)
//		db.DB().SetMaxOpenConns(setting.RunTime.Database.MaxOpenConns)
//		if err = Sync(); err != nil {
//			return err
//		}
//	}
//	return nil
//}

////Sync database structs
//func Sync() error {
//	log.Println("Sync database structs ")
//	db.AutoMigrate(&Pipeline{})
//	db.AutoMigrate(&PipelineVersion{})
//	db.AutoMigrate(&Point{})
//	db.AutoMigrate(&PointVersion{})
//	db.AutoMigrate(&Stage{})
//	db.AutoMigrate(&StageVersion{})
//	return nil
//}
