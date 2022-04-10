package template

const RPCServerTemplate = `package server

import (
	"context"
	"{{.Project}}/proto/{{.ServiceType}}/{{.Module}}/{{.PackageName}}"
	"{{.Project}}/{{.ServiceType}}/{{.Module}}/{{.PackageName}}/rpc/logic"
	"{{.Project}}/{{.ServiceType}}/{{.Module}}/{{.PackageName}}/rpc/svc"
)

type Server struct {
	svcCtx *svc.ServiceContext
    {{.PackageName}}.Unimplemented{{.ServiceName}}Server
}

func NewServer(svcCtx *svc.ServiceContext) *Server {
	return &Server{
		svcCtx: svcCtx,
	}
}
{{range $index, $method := .MethodSigns}}
func (s *Server) {{$method.MethodName}}{{.Sign}}{
	{{if eq $method.ISStream true}}
	l := logic.New{{$method.MethodName}}Logic(stream.Context(), s.svcCtx)
	{{else}}
	l := logic.New{{$method.MethodName}}Logic(ctx, s.svcCtx)
	{{end}}
	return l.{{$method.MethodName}}({{.Param}})
}
{{end}}

`

type MethodSign struct {
	MethodName string
	Sign       string
	Param      string
	ISStream   bool
}

type ServerTemplateParam struct {
	Project     string
	PackageName string
	ServiceName string
	Module      string
	ServiceType string
	MethodSigns []MethodSign
}
