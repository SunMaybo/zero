package template

const HttpSvcTemplate = `package svc

import (
	"{{.Project}}/apis/{{.Module}}/{{.ServiceName}}/config"
	"github.com/SunMaybo/zero/common/zrpc"
)

type ServiceContext struct {
	Cfg       *config.Config
	RpcClient *zrpc.Client
}

func NewServiceContext(cfg *config.Config) *ServiceContext {

	return &ServiceContext{
		Cfg:       cfg,
		RpcClient: zrpc.NewClient(cfg.Zero),
	}
}

`

type HttpSvcTemplateParam struct {
	Project     string
	Module      string
	ServiceName string
}
