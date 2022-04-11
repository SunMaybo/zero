package zrpc

import (
	"errors"
	"github.com/SunMaybo/zero/common/zcfg"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
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
	dataId   string
	OnChange func(namespace, group, dataId, data string)
}
type CenterClient interface {
	GetSchema() string
	DoRegister(instance ServiceInstance) error
	DeRegister(instance ServiceInstance) error
	Subscribe(param *SubscribeParam) error
	SelectInstances(instances SelectInstancesParam) ([]ServiceInstance, error)
	GetConfig(param ConfigParam) (string, error)
}
type ServerCenterClient struct {
	client *CenterClient
}

func NewSingleCenterClient(cfg *zcfg.SeverCenterConfig) (CenterClient, error) {
	if len(cfg.ServerConfigs) <= 0 {
		return nil, errors.New("server configs is empty")
	}
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
