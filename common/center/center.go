package center

import (
	"errors"
	"github.com/SunMaybo/zero/common/zcfg"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"os"
	"os/user"
	"runtime"
	"strings"
)

const (
	Nacos_Server_Center_Name = "nacos"
	Etcd_Server_Center_Name  = "etcd"
)

type ServiceInstance struct {
	ClusterName string
	ServiceName string
	Port        int
	GroupName   string
	Address     string
	Weight      float64
	Metadata    map[string]string
}
type SubscribeParam struct {
	ServiceName       string
	Clusters          []string
	GroupName         string
	SubscribeCallback func(services []ServiceInstance)
}

type SelectInstancesParam struct {
	Clusters    []string `param:"clusters"`
	ServiceName string   `param:"serviceName"`
	GroupName   string   `param:"groupName"`
	HealthyOnly bool     `param:"healthyOnly"`
}

type ConfigParam struct {
	Group    string
	DataId   string
	Refresh  bool
	OnChange func(group, dataId, data string)
}
type Center interface {
	GetSchema() string
	DoRegister(instance ServiceInstance) error
	DeRegister(instance ServiceInstance) error
	Subscribe(param *SubscribeParam) error
	SelectInstances(instances SelectInstancesParam) ([]ServiceInstance, error)
	GetConfig(param ConfigParam) (string, error)
}
type ServerCenterClient struct {
	client *Center
}

func NewSingleCenterClient(cfg zcfg.SeverCenterConfig) (Center, error) {
	if len(cfg.ServerConfigs) <= 0 {
		return nil, errors.New("server configs is empty")
	}
	u, _ := user.Current()
	if cfg.CacheDir == "" && u != nil {
		cfg.CacheDir = u.HomeDir + "/.nacos/cache"
	}
	if cfg.LogDir == "" && u != nil {
		cfg.LogDir = u.HomeDir + "/.nacos/log"
	}
	createDir(cfg.CacheDir)
	createDir(cfg.LogDir)
	switch cfg.ServerCenterName {
	case Nacos_Server_Center_Name:
		clientConfig := constant.ClientConfig{
			TimeoutMs:    cfg.TimeoutMs,
			BeatInterval: cfg.BeatInterval,
			NamespaceId:  cfg.NamespaceId,
			CacheDir:     cfg.CacheDir,
			Username:     cfg.Username,
			Password:     cfg.Password,
			LogDir:       cfg.LogDir,
			LogLevel:     cfg.LogLevel,
		}
		var serverConfigs []constant.ServerConfig
		for _, config := range cfg.ServerConfigs {
			serverConfigs = append(serverConfigs, constant.ServerConfig{
				Scheme:      config.Scheme,
				ContextPath: config.ContextPath,
				IpAddr:      config.IpAddr,
				Port:        config.Port,
			})
		}
		return NewSingleNacosClient(&clientConfig, serverConfigs)

	}
	return nil, errors.New("not support server center name:" + cfg.ServerCenterName)
}
func createDir(dir string) {
	if runtime.GOOS == "windows" {
		dir = strings.ReplaceAll(dir, "/", "\\")
	}
	_, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, 0755)
	}
}
