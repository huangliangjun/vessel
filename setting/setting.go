package setting

import (
	"io/ioutil"
	"log"

	"github.com/ghodss/yaml"
)

// GlobalConf global config
type GlobalConf struct {
	AppName     string
	Usage       string
	Version     string
	Author      string
	Email       string
	RuntimePath string
}

// RunTimeConf runtime config
type RunTimeConf struct {
	// Run config
	Run struct {
		runMode string
		logPath string
	}
	// Http config
	HTTP struct {
		ListenMode    string
		HTTPSCertFile string
		HTTPSKeyFile  string
		Host          string
		Port          string
	}
	// Database config
	Database struct {
		Username string
		Password string
		Protocol string
		Host     string
		Port     string
		Schema   string
		Param    map[string]string
	}
	// Etcd config
	Etcd struct {
		Endpoints []map[string]string
		Username  string
		Password  string
	}
	// K8s config
	K8s struct {
		Host string
		Port string
	}
}

var (
	// Global global config
	Global *GlobalConf
	// RunTime runTime config
	RunTime *RunTimeConf
)

// InitGlobalConf global config init
func InitGlobalConf(globalFilePath string) error {
	globalFile, err := ioutil.ReadFile(globalFilePath)
	if err != nil {
		return err
	}

	Global = &GlobalConf{}
	if err = yaml.Unmarshal([]byte(globalFile), &Global); err != nil {
		return err
	}

	return initRuntimeConf(Global.RuntimePath)
}

func initRuntimeConf(runtimeFilePath string) error {
	runtimeFile, err := ioutil.ReadFile(runtimeFilePath)
	if err != nil {
		return err
	}

	RunTime = &RunTimeConf{}
	if err := yaml.Unmarshal([]byte(runtimeFile), &RunTime); err != nil {
		return err
	}
	log.Println(RunTime)
	return nil
}
