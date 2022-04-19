package template

const RPCMainTemplate = `package main

import (
	"flag"
	"google.golang.org/grpc"
	"github.com/SunMaybo/zero/common/zcfg"
	"github.com/SunMaybo/zero/common/zrpc"
	"{{.Project}}/proto/{{.ServiceType}}/{{.Module}}/{{.PackageName}}"
	"{{.Project}}/{{.ServiceType}}/{{.Module}}/{{.PackageName}}/rpc/config"
	"{{.Project}}/{{.ServiceType}}/{{.Module}}/{{.PackageName}}/rpc/server"
	"{{.Project}}/{{.ServiceType}}/{{.Module}}/{{.PackageName}}/rpc/svc"
)

var cfgPath = flag.String("cfg", "./etc/config.yaml", "config file path")

func init() {
	flag.Parse()
}

func main() {
	cfg := config.Config{}
	zcfg.LoadConfig(*cfgPath, &cfg)
    //jwtInterceptor:=grpc.ChainUnaryInterceptor(
	//	interceptor.UnaryJWTServerInterceptor("",nil),
	//))
	s := zrpc.NewServer(cfg.Zero)
	defer s.Stop()
	s.RegisterServer(func(s *grpc.Server) error {
		serviceContext := svc.NewServiceContext(cfg)
		{{.PackageName}}.Register{{.Service}}Server(s, server.NewServer(serviceContext))
		return nil
	}).Start()
}

`

type MainTemplateParam struct {
	Project     string
	PackageName string
	Module      string
	ServiceType string
	Service     string
}
