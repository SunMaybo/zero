package template

const RPCMainTemplate = `package main

import (
	"flag"
	"google.golang.org/grpc"
	"github.com/SunMaybo/zero/common/zcfg"
	"github.com/SunMaybo/zero/common/zrpc"
	"{{.Project}}/proto/{{.Module}}/{{.PackageName}}"
	"{{.Project}}/{{.ServiceType}}/{{.Module}}/{{.PackageName}}/rpc/config"
	"{{.Project}}/{{.ServiceType}}/{{.Module}}/{{.PackageName}}/rpc/server"
	"{{.Project}}/{{.ServiceType}}/{{.Module}}/{{.PackageName}}/rpc/svc"
)

var cfgPath = flag.String("cfg", "etc/config.yaml", "cfg path")

func main() {
	flag.Parse()
	cfg := zcfg.LoadConfig[config.Config](*cfgPath)
	s := zrpc.NewServer(cfg.Rpc, []grpc.UnaryServerInterceptor{
		//zrpc.UnaryJWTServerInterceptor("", []zrpc.RpcMethod{
		//	greeter.Rpc_SayStream,
		//}),
	}, nil)
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
