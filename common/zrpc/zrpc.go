package zrpc

import (
	"context"
	"fmt"
	"github.com/SunMaybo/zero/common/zlog"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_revovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"
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

type HystrixConfigTable map[string]*HystrixConfig
type RpcCfg struct {
	SeverCenterConfig SeverCenterConfig  `yaml:"center"`
	Hystrix           HystrixConfigTable `yaml:"hystrix"`
	Name              string             `yaml:"name"`
	Port              int                `yaml:"port"`
	Weight            float64            `yaml:"weight"`
	IsOnline          bool               `yaml:"is_online"`
	Metadata          map[string]string  `yaml:"metadata"`
	ClusterName       string             `yaml:"cluster_name"` // the cluster name
	GroupName         string             `yaml:"group_name"`   // the group name
	Timeout           int                `yaml:"timeout"`
	EnableMetrics     bool               `yaml:"enable_metrics"`
	MetricsPort       int                `yaml:"metrics_port"`
	MetricsPath       string             `yaml:"metrics_path"`
}

type Server struct {
	grpcServer *grpc.Server
	logger     *zap.Logger
	isRegister *atomic.Bool
	rpcCfg     RpcCfg
	tracer     opentracing.Tracer
}

func NewServer(cfg RpcCfg, options ...grpc.ServerOption) *Server {
	tracer, _ := jaeger.NewTracer(
		"grpc",
		jaeger.NewConstSampler(true),
		jaeger.NewNullReporter(),
	)
	return NewServerWithTracer(cfg, tracer, options...)
}

func NewServerWithTracer(cfg RpcCfg, tracer opentracing.Tracer, options ...grpc.ServerOption) *Server {
	// init logger
	zlog.InitLogger(cfg.IsOnline)
	// setting grpc server timeout
	if cfg.Timeout <= 0 {
		cfg.Timeout = 5000
	}
	var defaultOptions = []grpc.UnaryServerInterceptor{
		NewValidatorInterceptor().Interceptor,
		grpc_revovery.UnaryServerInterceptor(),
		otgrpc.OpenTracingServerInterceptor(tracer),
		UnaryLoggerServerInterceptor(),
		UnaryTimeoutInterceptor(time.Duration(cfg.Timeout) * time.Millisecond),
	}
	defaultStreamOptions := []grpc.StreamServerInterceptor{
		grpc_revovery.StreamServerInterceptor(),
		otgrpc.OpenTracingStreamServerInterceptor(tracer),
		StreamLoggerServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
	}
	if cfg.EnableMetrics {
		defaultOptions = append(defaultOptions, grpc_prometheus.UnaryServerInterceptor)
		defaultStreamOptions = append(defaultStreamOptions, grpc_prometheus.StreamServerInterceptor)
		//begin prometheus metrics
		go bindingMetrics(cfg.MetricsPath, cfg.MetricsPort)
	}
	options = append(options, grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(defaultOptions...),
	), grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(defaultStreamOptions...),
	))
	return &Server{
		grpcServer: grpc.NewServer(options...),
		rpcCfg:     cfg,
		isRegister: atomic.NewBool(false),
		logger:     zlog.LOGGER,
		tracer:     tracer,
	}
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
				logger.Info("receive signal", zap.String("signal", sign.String()))
				if s.isRegister.Load() {
					s.unRegister()
				}
				//退出
				os.Exit(0)
			default:
			}
		}
	}()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.rpcCfg.Port))
	if err != nil {
		s.logger.Fatal("failed to listen", zap.Error(err))
		return
	}
	s.logger.Info("start to serve", zap.String("name", s.rpcCfg.Name), zap.Int("port", s.rpcCfg.Port))
	if err := s.grpcServer.Serve(lis); err != nil {
		s.logger.Fatal("failed to serve", zap.Error(err))
		return
	}
}
func (s *Server) register() {
	if !s.rpcCfg.SeverCenterConfig.Enable {
		return
	}
	time.Sleep(3 * time.Second)
	center, err := NewSingleCenterClient(s.rpcCfg.SeverCenterConfig)
	if err != nil {
		panic(err)
	}
	NewRegister(center).DoRegister(ServiceInstance{
		ServiceName: s.rpcCfg.Name,
		Port:        s.rpcCfg.Port,
		Weight:      s.rpcCfg.Weight,
		ClusterName: s.rpcCfg.ClusterName,
		GroupName:   s.rpcCfg.GroupName,
		Metadata:    s.rpcCfg.Metadata,
	})
	s.isRegister.Store(true)
	s.logger.Info("register success", zap.String("name", s.rpcCfg.Name), zap.Int("port", s.rpcCfg.Port))
}
func (s *Server) unRegister() {
	if !s.rpcCfg.SeverCenterConfig.Enable {
		return
	}
	center, err := NewSingleCenterClient(s.rpcCfg.SeverCenterConfig)
	if err != nil {
		panic(err)
	}
	NewRegister(center).UnRegister(ServiceInstance{
		ServiceName: s.rpcCfg.Name,
		Port:        s.rpcCfg.Port,
		Weight:      s.rpcCfg.Weight,
		ClusterName: s.rpcCfg.ClusterName,
		GroupName:   s.rpcCfg.GroupName,
		Metadata:    s.rpcCfg.Metadata,
	})
	s.isRegister.Store(false)
	s.logger.Info("unregister success", zap.String("name", s.rpcCfg.Name), zap.Int("port", s.rpcCfg.Port))
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
	hystrix      HystrixConfigTable
	tracer       opentracing.Tracer
}

func NewClient(cfg RpcCfg) *Client {
	tracer, _ := jaeger.NewTracer(
		"grpc",
		jaeger.NewConstSampler(true),
		jaeger.NewNullReporter(),
	)
	return NewClientWithTracer(cfg, tracer)

}
func NewClientWithTracer(cfg RpcCfg, tracer opentracing.Tracer) *Client {
	zlog.InitLogger(cfg.IsOnline)
	center, err := NewSingleCenterClient(cfg.SeverCenterConfig)
	if err != nil {
		zlog.S.Errorf("connection discovery center failed,err:%s", err.Error())
	}
	resolver.Register(NewResolverBuilder(center))
	return &Client{
		clusterNames: cfg.ClusterName,
		groupName:    cfg.GroupName,
		schema:       cfg.SeverCenterConfig.ServerCenterName,
		hystrix:      cfg.Hystrix,
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
			UnaryLoggerClientInterceptor(),
			TimeoutInterceptor(timeout),
			UnaryHystrixClientInterceptor(hystrixCfg),
		),
		grpc.WithChainStreamInterceptor(
			otgrpc.OpenTracingStreamClientInterceptor(c.tracer),
			StreamLoggerClientInterceptor(),
		),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	return grpc.DialContext(ctx, serviceName, options...)
}
