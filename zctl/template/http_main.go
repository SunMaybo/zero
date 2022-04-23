package template

const HttpMainTemplate = `package main

import (
	"flag"
	"{{.Project}}/apis/{{.Module}}/{{.ServiceName}}/config"
	"{{.Project}}/apis/{{.Module}}/{{.ServiceName}}/server"
	"{{.Project}}/apis/{{.Module}}/{{.ServiceName}}/svc"
	"github.com/SunMaybo/zero/common/zcfg"
)

// 设置路由信息
var path string

func init() {
	flag.StringVar(&path, "etc", "etc/config.yaml", "config for filepath")
	flag.Parse()
}
func main() {
	cfg := config.Config{}
	zcfg.LoadConfig(path, &cfg)
	svcCtx := svc.NewServiceContext(&cfg)
	server.NewServer(svcCtx).Start()
}

`

type HttpMainTemplateParam struct {
	Project     string
	Module      string
	ServiceName string
}
