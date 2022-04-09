package template

const RPCSvcTemplate = `package svc

import "{{.Project}}/{{.ServiceType}}/{{.Module}}/{{.PackageName}}/rpc/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(cfg config.Config) *ServiceContext {
	return &ServiceContext{
		Config: cfg,
	}
}
`

type SvcTemplateParam struct {
	Project     string
	PackageName string
	Module      string
	ServiceType string
}
