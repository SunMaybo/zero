package zgin

import (
	"context"
	"fmt"
	"github.com/SunMaybo/zero/common/zlog"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	"net/http"
	"time"
)

type Server struct {
	engine  *gin.Engine
	port    int64
	timeout time.Duration
}

func NewServer(port int64) *Server {
	t, err := zipkin.NewTracer(nil)
	if err != nil {
		zlog.S.Fatal(err)
	}
	t.SetNoop(false)
	tracer := zipkinot.Wrap(t)
	opentracing.SetGlobalTracer(tracer)
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(zlog.RecoveryWithZap())
	r.Use(Cors())
	return &Server{
		engine:  r,
		port:    port,
		timeout: time.Second * 30,
	}

}
func (s *Server) SetTracer(tracer opentracing.Tracer) {
	opentracing.SetGlobalTracer(tracer)
}
func NewServerWithTimeout(port int64, timeout time.Duration, mode string) *Server {
	gin.SetMode(mode)
	r := gin.New()
	r.Use(zlog.RecoveryWithZap())
	r.Use(Cors())
	return &Server{
		engine:  r,
		port:    port,
		timeout: timeout,
	}

}

type ZeroGinHandler func(ctx context.Context, ginCtx *gin.Context)

func (s *Server) MiddleHandle(handler ZeroGinHandler) func(context *gin.Context) {
	return func(ginCtx *gin.Context) {
		var ParentSpan opentracing.Span
		spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(ginCtx.Request.Header))
		if err != nil {
			ParentSpan = opentracing.GlobalTracer().StartSpan(ginCtx.Request.URL.Path)
			defer ParentSpan.Finish()
		} else {
			ParentSpan = opentracing.StartSpan(
				ginCtx.Request.URL.Path,
				opentracing.ChildOf(spCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				ext.SpanKindRPCServer,
			)
			defer ParentSpan.Finish()
		}
		ctx, cancel := context.WithTimeout(context.TODO(), s.timeout)
		defer cancel()
		handler(opentracing.ContextWithSpan(ctx, ParentSpan), ginCtx)
	}
}

func (s *Server) Start(router func(engine *gin.Engine)) {
	router(s.engine)
	zlog.S.Infof("start server at %d", s.port)
	if err := s.engine.Run(fmt.Sprintf(":%d", s.port)); err != nil {
		zlog.S.Fatal(err)
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
