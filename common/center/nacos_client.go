package center

import (
	"errors"
	"fmt"
	"github.com/SunMaybo/zero/common/zlog"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"sync"
)

var (
	nacosClient     *NacosClient
	onceNacosClient = sync.Once{}
)

type NacosClient struct {
	nameClient   naming_client.INamingClient
	configClient config_client.IConfigClient
}

func NewSingleNacosClient(clientConfig *constant.ClientConfig, serverConfigs []constant.ServerConfig) (*NacosClient, error) {
	var err error
	onceNacosClient.Do(
		func() {
			var nameClient naming_client.INamingClient
			var configClient config_client.IConfigClient
			nameClient, err = clients.NewNamingClient(vo.NacosClientParam{
				ClientConfig:  clientConfig,
				ServerConfigs: serverConfigs,
			})
			if err != nil {
				return
			}
			configClient, err = clients.NewConfigClient(vo.NacosClientParam{
				ClientConfig:  clientConfig,
				ServerConfigs: serverConfigs,
			})
			if err != nil {
				return
			}
			nacosClient = &NacosClient{
				nameClient:   nameClient,
				configClient: configClient,
			}
		})
	if err != nil || nacosClient == nil {
		return nil, errors.New(fmt.Sprintf("create naming client failed, err: %s", err.Error()))
	}
	return nacosClient, nil
}
func (n *NacosClient) DoRegister(instance ServiceInstance) error {
	if isOk, err := n.nameClient.RegisterInstance(vo.RegisterInstanceParam{
		ServiceName: instance.ServiceName,
		Ip:          instance.Address,
		Port:        uint64(instance.Port),
		Weight:      instance.Weight,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    instance.Metadata,
		ClusterName: instance.ClusterName,
		GroupName:   instance.GroupName,
	}); err != nil {
		return err
	} else if !isOk {
		return errors.New("register instance failed")
	}
	return nil
}
func (n *NacosClient) DeRegister(instance ServiceInstance) error {
	if isOk, err := n.nameClient.DeregisterInstance(vo.DeregisterInstanceParam{
		ServiceName: instance.ServiceName,
		Ip:          instance.Address,
		Port:        uint64(instance.Port),
		GroupName:   instance.GroupName,
		Cluster:     instance.ClusterName,
		Ephemeral:   true,
	}); err != nil {
		return err
	} else if !isOk {
		return errors.New("deregister instance failed")
	}
	return nil
}
func (n *NacosClient) Subscribe(param *SubscribeParam) error {
	return n.nameClient.Subscribe(&vo.SubscribeParam{
		ServiceName: param.ServiceName,
		Clusters:    param.Clusters,
		GroupName:   param.GroupName,
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			if err != nil {
				return
			}
			var serviceInstances []ServiceInstance
			for _, service := range services {
				if service.Healthy && service.Enable {
					serviceInstances = append(serviceInstances, ServiceInstance{
						ServiceName: service.ServiceName,
						Address:     service.Ip,
						Port:        int(service.Port),
						Weight:      service.Weight,
						Metadata:    service.Metadata,
						ClusterName: service.ClusterName,
					})
				}
			}
			param.SubscribeCallback(serviceInstances)
		},
	})
}
func (n *NacosClient) SelectInstances(instance SelectInstancesParam) ([]ServiceInstance, error) {
	instances, err := n.nameClient.SelectInstances(vo.SelectInstancesParam{
		Clusters:    instance.Clusters,
		GroupName:   instance.GroupName,
		ServiceName: instance.ServiceName,
		HealthyOnly: instance.HealthyOnly,
	})
	if err != nil {
		return nil, err
	}
	var serviceInstances []ServiceInstance
	for _, m := range instances {
		if m.Healthy && m.Enable {
			serviceInstances = append(serviceInstances, ServiceInstance{
				ServiceName: m.ServiceName,
				Address:     m.Ip,
				Port:        int(m.Port),
				Weight:      m.Weight,
				Metadata:    m.Metadata,
				ClusterName: m.ClusterName,
			})
		}
	}
	return serviceInstances, err
}

func (n *NacosClient) GetConfig(param ConfigParam) (string, error) {
	data, err := n.configClient.GetConfig(vo.ConfigParam{
		Group:  param.Group,
		DataId: param.DataId,
		Type:   vo.YAML,
	})
	if data != "" {
		param.OnChange(param.Group, param.DataId, data)
	}
	if err != nil {
		zlog.S.Warnw("get config failed", "err", err)
	}
	if param.Refresh {
		err = n.configClient.ListenConfig(vo.ConfigParam{
			Group:  param.Group,
			DataId: param.DataId,
			Type:   vo.YAML,
			OnChange: func(namespace, group, dataId, data string) {
				param.OnChange(param.Group, param.DataId, data)
			},
		})
	}
	return data, err
}

func (n *NacosClient) GetSchema() string {
	return "nacos"
}
