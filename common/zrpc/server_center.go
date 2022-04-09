package zrpc

import "github.com/nacos-group/nacos-sdk-go/common/constant"

const (
	Nacos_Server_Center_Name = "nacos"
	Etcd_Server_Center_Name  = "etcd"
)

type SeverCenterConfig struct {
	TimeoutMs        uint64         `yaml:"timeout_ms"`    // timeout for requesting Nacos server, default value is 10000ms
	BeatInterval     int64          `yaml:"beat_interval"` // the time interval for sending beat to server,default value is 5000ms
	NamespaceId      string         `yaml:"namespace_id"`  // the namespaceId of Nacos.When namespace is public, fill in the blank string here.
	CacheDir         string         `yaml:"cache_dir"`     // the directory for persist nacos service info,default value is current path
	Username         string         `yaml:"username"`      // the username for nacos auth
	Password         string         `yaml:"password"`      // the password for nacos auth
	LogDir           string         `yaml:"log_dir"`       // the directory for log, default is current path
	LogLevel         string         `yaml:"log_level"`     // the level of log, it's must be debug,info,warn,error, default value is info
	Enable           bool           `yaml:"enabled"`       // enable or disable the server center
	ServerConfigs    []ServerConfig `yaml:"server"`        // the server configs
	ServerCenterName string         `yaml:"name"`          // the server center name, default value is Nacos_Server_Center
}
type ServerConfig struct {
	Scheme      string `yaml:"scheme"`       //the server scheme
	ContextPath string `yaml:"context_path"` //the server contextpath
	IpAddr      string `yaml:"host"`         //the server address
	Port        uint64 `yaml:"port"`         //the server port
}
type ServiceInstance struct {
	ClusterName string
	ServiceName string
	Port        int
	GroupName   string
	Address     string
	Weight      float64
	Metadata    map[string]string
}
type CenterClient interface {
	DoRegister(instance ServiceInstance) error
	DeRegister(instance ServiceInstance) error
}
type ServerCenterClient struct {
	client *CenterClient
}

func NewCenterClient(cfg SeverCenterConfig) CenterClient {
	if cfg.Enable {
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
			return NewNacosClient(&clientConfig, serverConfigs)

		}
	}
	return nil
}
