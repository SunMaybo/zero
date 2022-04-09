package zrpc

import (
	"fmt"
	"github.com/SunMaybo/zero/common/zlog"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_revovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type RpcCfg struct {
	SeverCenterConfig SeverCenterConfig `yaml:"center"`
	Name              string            `yaml:"name"`
	Port              int               `yaml:"port"`
	Weight            float64           `yaml:"weight"`
	IsOnline          bool              `yaml:"is_online"`
	Metadata          map[string]string `yaml:"metadata"`
	ClusterName       string            `yaml:"cluster_name"` // the cluster name
	GroupName         string            `yaml:"group_name"`   // the group name
	Timeout           int               `yaml:"timeout"`
	EnableMetrics     bool              `yaml:"enable_metrics"`
	MetricsPort       int               `yaml:"metrics_port"`
	MetricsPath       string            `yaml:"metrics_path"`
}

type Server struct {
	grpcServer *grpc.Server
	logger     *zap.Logger
	center     CenterClient
	isRegister *atomic.Bool
	rpcCfg     RpcCfg
}

var onceCenter = sync.Once{}

func NewServer(cfg RpcCfg, unaryServerInterceptors []grpc.UnaryServerInterceptor, streamServerInterceptors []grpc.StreamServerInterceptor) *Server {
	// init logger
	zlog.InitLogger(cfg.IsOnline)
	// setting grpc server timeout
	if cfg.Timeout <= 0 {
		cfg.Timeout = 5000
	}
	var options = []grpc.UnaryServerInterceptor{
		grpc_zap.UnaryServerInterceptor(zlog.LOGGER),
		NewValidatorInterceptor().Interceptor,
		grpc_revovery.UnaryServerInterceptor(),
		UnaryTimeoutInterceptor(time.Duration(cfg.Timeout) * time.Millisecond),
	}
	options = append(options, unaryServerInterceptors...)

	streamOptions := []grpc.StreamServerInterceptor{
		grpc_zap.StreamServerInterceptor(zlog.LOGGER),
		grpc_revovery.StreamServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
	}
	streamOptions = append(streamOptions, streamServerInterceptors...)

	if cfg.EnableMetrics {
		options = append(options, grpc_prometheus.UnaryServerInterceptor)
		streamOptions = append(streamOptions, grpc_prometheus.StreamServerInterceptor)
		//begin prometheus metrics
		go bindingMetrics(cfg.MetricsPath, cfg.MetricsPort)
	}

	return &Server{
		grpcServer: grpc.NewServer(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(options...)),
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamOptions...)),
		),
		rpcCfg:     cfg,
		isRegister: atomic.NewBool(false),
		logger:     zlog.LOGGER,
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
	signal.Notify(signChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for sign := range signChan {
			switch sign {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				logger.Info("receive signal", zap.String("signal", sign.String()))
				if s.isRegister.Load() {
					s.unRegister()
				}
				//退出
				os.Exit(0)
			case syscall.SIGUSR1:
			case syscall.SIGUSR2:
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
	onceCenter.Do(func() {
		s.center = NewCenterClient(s.rpcCfg.SeverCenterConfig)
	})
	NewRegister(s.center).DoRegister(ServiceInstance{
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
	onceCenter.Do(func() {
		s.center = NewCenterClient(s.rpcCfg.SeverCenterConfig)
	})
	NewRegister(s.center).UnRegister(ServiceInstance{
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
