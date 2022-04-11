package zcfg

type ZeroConfig struct {
	RPC               RpcCfg            `yaml:"rpc"`
	SeverCenterConfig SeverCenterConfig `yaml:"center"`
}
type HystrixConfigTable map[string]*HystrixConfig

type HystrixConfig struct {
	Timeout                int `yaml:"timeout"`
	MaxConcurrentRequests  int `yaml:"max_concurrent_requests"`
	SleepWindow            int `yaml:"sleep_window"`
	RequestVolumeThreshold int `yaml:"request_volume_threshold"`
	MaxRetry               int `yaml:"max_retry"`
	RetryTimeout           int `yaml:"retry_timeout"`
	ErrorPercentThreshold  int `yaml:"error_percent_threshold"`
}

type RpcCfg struct {
	Hystrix       HystrixConfigTable `yaml:"hystrix"`
	Name          string             `yaml:"name"`
	Port          int                `yaml:"port"`
	Weight        float64            `yaml:"weight"`
	IsOnline      bool               `yaml:"is_online"`
	Metadata      map[string]string  `yaml:"metadata"`
	ClusterName   string             `yaml:"cluster_name"` // the cluster name
	GroupName     string             `yaml:"group_name"`   // the group name
	Timeout       int                `yaml:"timeout"`
	EnableMetrics bool               `yaml:"enable_metrics"`
	MetricsPort   int                `yaml:"metrics_port"`
	MetricsPath   string             `yaml:"metrics_path"`
}
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
