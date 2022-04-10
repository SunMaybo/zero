package template

const RPCLogicTemplate = `package logic

import (
	"context"
	"go.uber.org/zap"
	"github.com/SunMaybo/zero/common/zlog"
	"{{.Project}}/proto/{{.ServiceType}}/{{.Module}}/{{.PackageName}}"
	"{{.Project}}/{{.ServiceType}}/{{.Module}}/{{.PackageName}}/rpc/svc"
)

type {{.MethodName}}Logic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *zap.SugaredLogger
}

func New{{.MethodName}}Logic(ctx context.Context, svcCtx *svc.ServiceContext) *{{.MethodName}}Logic {
	return &{{.MethodName}}Logic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: zlog.WithContext(ctx),
	}
}
func (l *{{.MethodName}}Logic) {{.MethodName}}{{.Sign}} {
	return {{.Return}}
}
`

type LogicTemplateParam struct {
	Project     string
	PackageName string
	ServiceType string
	MethodName  string
	Sign        string
	Return      string
	Module      string
}
