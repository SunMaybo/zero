package template

const HttpServerTemplate = `package server

import (
	"context"
	"{{.Project}}/apis/{{.Module}}/{{.ServiceName}}/svc"
	"github.com/SunMaybo/zero/common/zgin"
	"github.com/gin-gonic/gin"
	"time"
)

type Server struct {
	svcCtx *svc.ServiceContext
	server *zgin.Server
}

func NewServer(svc *svc.ServiceContext) *Server {
	if svc.Cfg.Zero.Server.Timeout <= 0 {
		svc.Cfg.Zero.Server.Timeout = 30
	}
	if svc.Cfg.Zero.Server.Port <= 0 {
		svc.Cfg.Zero.Server.Port = 3000
	}
	return &Server{
		svcCtx: svc,
		server: zgin.NewServerWithTimeout(svc.Cfg.Zero.Server.Port, time.Duration(svc.Cfg.Zero.Server.Timeout)*time.Second),
	}
}
func (s *Server) Start() {
	s.server.Start(func(engine *gin.Engine) {
		engine.GET("/ping", s.server.MiddleHandle(func(ctx context.Context, ginCtx *gin.Context) {
			ginCtx.JSON(200, gin.H{
				"message": "pong",
			})
		}))
	})
}

`

type HttpServerTemplateParam struct {
	Project     string
	ServiceName string
	Module      string
}
