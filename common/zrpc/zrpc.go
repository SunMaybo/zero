package zrpc

import (
	"context"
	"fmt"
	"github.com/SunMaybo/zero/common/center"
	"github.com/SunMaybo/zero/common/zcfg"
	"github.com/SunMaybo/zero/common/zlog"
	"github.com/SunMaybo/zero/common/zrpc/interceptor"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_revovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/resolver"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	grpcServer   *grpc.Server
	logger       *zap.Logger
	isRegister   *atomic.Bool
	zeroCfg      zcfg.ZeroConfig
	tracer       opentracing.Tracer
	center       center.Center
	configParams []center.ConfigParam
}

func NewServer(cfg zcfg.ZeroConfig, options ...grpc.ServerOption) *Server {
	tarcer, _ := zipkin.NewTracer(nil)
	return NewServerWithTracer(cfg, zipkinot.Wrap(tarcer), options...)
}

func NewServerWithTracer(cfg zcfg.ZeroConfig, tracer opentracing.Tracer, options ...grpc.ServerOption) *Server {
	// init logger
	zlog.InitLogger(cfg.RPC.IsOnline)
	// setting grpc server timeout
	if cfg.RPC.Timeout <= 0 {
		cfg.RPC.Timeout = 5000
	}
	var defaultOptions = []grpc.UnaryServerInterceptor{
		otgrpc.OpenTracingServerInterceptor(tracer),
		grpc_revovery.UnaryServerInterceptor(grpc_revovery.WithRecoveryHandlerContext(func(ctx context.Context, p interface{}) (err error) {
			zlog.WithContext(ctx).Errorf("err:%s", p)
			return nil
		})),
		interceptor.NewValidatorInterceptor().Interceptor,
		interceptor.UnaryLoggerServerInterceptor(),
		interceptor.UnaryTimeoutInterceptor(time.Duration(cfg.RPC.Timeout) * time.Millisecond),
	}
	defaultStreamOptions := []grpc.StreamServerInterceptor{
		otgrpc.OpenTracingStreamServerInterceptor(tracer),
		grpc_revovery.StreamServerInterceptor(grpc_revovery.WithRecoveryHandlerContext(func(ctx context.Context, p interface{}) (err error) {
			zlog.WithContext(ctx).Errorf("err:%s", p)
			return nil
		})),
		grpc_revovery.StreamServerInterceptor(),
		interceptor.StreamLoggerServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
	}
	if cfg.RPC.EnableMetrics {
		defaultOptions = append(defaultOptions, grpc_prometheus.UnaryServerInterceptor)
		defaultStreamOptions = append(defaultStreamOptions, grpc_prometheus.StreamServerInterceptor)
		//begin prometheus metrics
		go bindingMetrics(cfg.RPC.MetricsPath, cfg.RPC.MetricsPort)
	}
	center, err := center.NewSingleCenterClient(&cfg.SeverCenterConfig)
	if err != nil {
		zlog.S.Warnw("create center client failed", "err", err)
	}
	var allOptions []grpc.ServerOption
	allOptions = append(allOptions, grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(defaultOptions...),
	), grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(defaultStreamOptions...),
	))
	allOptions = append(allOptions, options...)
	return &Server{
		grpcServer: grpc.NewServer(allOptions...),
		zeroCfg:    cfg,
		isRegister: atomic.NewBool(false),
		logger:     zlog.LOGGER,
		tracer:     tracer,
		center:     center,
	}
}

func (s *Server) AddConfigListener(configParam ...center.ConfigParam) {
	s.configParams = append(s.configParams, configParam...)
}

type RegisterFunc func(s *grpc.Server) error

func (s *Server) RegisterServer(registerFunc RegisterFunc) *Server {
	if err := registerFunc(s.grpcServer); err != nil {
		s.logger.Fatal("failed to register server", zap.Error(err))
	}
	reflection.Register(s.grpcServer)
	go s.register()
	return s
}
func (s *Server) Start() {
	//创建监听退出chan
	signChan := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(signChan, os.Interrupt, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for sign := range signChan {
			switch sign {
			case os.Interrupt, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR2, syscall.SIGUSR1:
				s.logger.Info("receive signal", zap.String("signal", sign.String()))
				if s.isRegister.Load() {
					s.unRegister()
				}
				//退出
				os.Exit(0)
			default:
			}
		}
	}()
	//监听配置
	for _, param := range s.configParams {
		if s.center != nil {
			if _, err := s.center.GetConfig(param); err != nil {
				s.logger.Sugar().Warnw("get config failed", "group", param.Group, "data_id", param.DataId, "err", err)
			}
		}
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.zeroCfg.RPC.Port))
	if err != nil {
		s.logger.Fatal("failed to listen", zap.Error(err))
		return
	}
	s.logger.Info("start to serve", zap.String("name", s.zeroCfg.RPC.Name), zap.Int("port", s.zeroCfg.RPC.Port))
	if err := s.grpcServer.Serve(lis); err != nil {
		s.logger.Fatal("failed to serve", zap.Error(err))
		return
	}
}
func (s *Server) register() {
	if !s.zeroCfg.SeverCenterConfig.Enable {
		return
	}
	time.Sleep(3 * time.Second)
	NewRegister(s.center).DoRegister(center.ServiceInstance{
		ServiceName: s.zeroCfg.RPC.Name,
		Port:        s.zeroCfg.RPC.Port,
		Weight:      s.zeroCfg.RPC.Weight,
		ClusterName: s.zeroCfg.RPC.ClusterName,
		GroupName:   s.zeroCfg.RPC.GroupName,
		Metadata:    s.zeroCfg.RPC.Metadata,
	})
	s.isRegister.Store(true)
	s.logger.Info("register success", zap.String("name", s.zeroCfg.RPC.Name), zap.Int("port", s.zeroCfg.RPC.Port))
}
func (s *Server) unRegister() {
	if !s.zeroCfg.SeverCenterConfig.Enable {
		return
	}
	NewRegister(s.center).UnRegister(center.ServiceInstance{
		ServiceName: s.zeroCfg.RPC.Name,
		Port:        s.zeroCfg.RPC.Port,
		Weight:      s.zeroCfg.RPC.Weight,
		ClusterName: s.zeroCfg.RPC.ClusterName,
		GroupName:   s.zeroCfg.RPC.GroupName,
		Metadata:    s.zeroCfg.RPC.Metadata,
	})
	s.isRegister.Store(false)
	s.logger.Info("unregister success", zap.String("name", s.zeroCfg.RPC.Name), zap.Int("port", s.zeroCfg.RPC.Port))
}
func bindingMetrics(metricPath string, metricPort int) {
	if metricPath == "" {
		metricPath = "/metrics"
	}
	if metricPort <= 0 {
		metricPort = 8848
	}
	http.Handle(metricPath, promhttp.Handler())
	_ = http.ListenAndServe(fmt.Sprintf(":%d", metricPort), nil)
}
func (s *Server) Stop() {
	s.grpcServer.Stop()
}

type Client struct {
	clusterNames string
	groupName    string
	schema       string
	hystrix      zcfg.HystrixConfigTable
	tracer       opentracing.Tracer
}

func NewClient(cfg zcfg.ZeroConfig) *Client {
	tarcer, _ := zipkin.NewTracer(nil)
	return NewClientWithTracer(cfg, zipkinot.Wrap(tarcer))

}
func NewClientWithTracer(cfg zcfg.ZeroConfig, tracer opentracing.Tracer) *Client {
	zlog.InitLogger(cfg.RPC.IsOnline)
	center, err := center.NewSingleCenterClient(&cfg.SeverCenterConfig)
	if err != nil {
		zlog.S.Errorf("connection discovery center failed,err:%s", err.Error())
	}
	resolver.Register(NewResolverBuilder(center))
	return &Client{
		clusterNames: cfg.RPC.ClusterName,
		groupName:    cfg.RPC.GroupName,
		schema:       cfg.SeverCenterConfig.ServerCenterName,
		hystrix:      cfg.RPC.Hystrix,
		tracer:       tracer,
	}

}
func (c *Client) GetGrpcClient(serviceName string, options ...grpc.DialOption) (*grpc.ClientConn, error) {
	return c.GetGrpcClientWithTimeout(serviceName, 3*time.Second, options...)
}
func (c *Client) GetGrpcClientWithTimeout(serviceName string, timeout time.Duration, options ...grpc.DialOption) (*grpc.ClientConn, error) {
	hystrixCfg := c.hystrix[serviceName]
	if !strings.HasPrefix(serviceName, c.schema+"://") {
		serviceName = fmt.Sprintf(c.schema+"://"+serviceName+"?cluster=%s&group=%s", c.clusterNames, c.groupName)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	options = append(options,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(c.tracer),
			interceptor.UnaryLoggerClientInterceptor(),
			interceptor.TimeoutInterceptor(timeout),
			interceptor.UnaryHystrixClientInterceptor(hystrixCfg),
		),
		grpc.WithChainStreamInterceptor(
			otgrpc.OpenTracingStreamClientInterceptor(c.tracer),
			interceptor.StreamLoggerClientInterceptor(),
		),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	return grpc.DialContext(ctx, serviceName, options...)
}
