package zrpc

import (
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type NacosClient struct {
	nameClient naming_client.INamingClient
}

func NewNacosClient(clientConfig *constant.ClientConfig, serverConfigs []constant.ServerConfig) *NacosClient {
	nameClient, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		panic(errors.New(fmt.Sprintf("create naming client failed, err: %s", err.Error())))
	}
	return &NacosClient{
		nameClient: nameClient,
	}
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
