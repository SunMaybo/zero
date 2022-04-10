package zrpc

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type HystrixConfig struct {
	Timeout                int `yaml:"timeout"`
	MaxConcurrentRequests  int `yaml:"max_concurrent_requests"`
	SleepWindow            int `yaml:"sleep_window"`
	RequestVolumeThreshold int `yaml:"request_volume_threshold"`
	MaxRetry               int `yaml:"max_retry"`
	RetryTimeout           int `yaml:"retry_timeout"`
	ErrorPercentThreshold  int `yaml:"error_percent_threshold"`
}

var lockMethodHystrix = sync.Mutex{}

var methodHystrixConfigTable = map[string]struct{}{}

func initHystrixConfigCommand(method string, cfg *HystrixConfig) {
	if _, ok := methodHystrixConfigTable[method]; ok {
		return
	}
	lockMethodHystrix.Lock()
	defer lockMethodHystrix.Unlock()
	if _, ok := methodHystrixConfigTable[method]; ok {
		return
	}
	if cfg != nil {
		hystrix.ConfigureCommand(method, hystrix.CommandConfig{
			Timeout:               cfg.Timeout,
			MaxConcurrentRequests: cfg.MaxConcurrentRequests,

			SleepWindow: cfg.SleepWindow,

			RequestVolumeThreshold: cfg.RequestVolumeThreshold,

			ErrorPercentThreshold: cfg.ErrorPercentThreshold,
		})
	} else {
		hystrix.ConfigureCommand(method, hystrix.CommandConfig{})
	}
	methodHystrixConfigTable[method] = struct{}{}

}

func UnaryHystrixClientInterceptor(cfg *HystrixConfig) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		maxRetry := 3
		retryTimeout := 50 * time.Millisecond
		if cfg != nil {
			if cfg.MaxRetry <= 0 {
				maxRetry = 1
			} else {
				maxRetry = cfg.MaxRetry
			}
			if cfg.RetryTimeout > 0 {
				retryTimeout = time.Duration(cfg.RetryTimeout) * time.Second
			}
		}
		initHystrixConfigCommand(method, cfg)
		var err error
		for i := 0; i < maxRetry; i++ {
			if err = hystrix.Do(method, func() error {
				return invoker(ctx, method, req, reply, cc, opts...)

			}, func(err error) error {
				return err
			}); err != nil {
				time.Sleep(retryTimeout)
				continue
			}
			break
		}
		return err
	}
}
