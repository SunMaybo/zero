package zrpc

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type HystrixConfig struct {
	Timeout                int
	MaxConcurrentRequests  int
	SleepWindow            int
	RequestVolumeThreshold int
	MaxRetry               int
	RetryTimeout           time.Duration
	ErrorPercentThreshold  int
}
type HystrixConfigTable map[string]HystrixConfig

var once sync.Once

func UnaryHystrixClientInterceptor(table HystrixConfigTable) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		maxRetry := 3
		retryTimeout := 50 * time.Millisecond
		if cfg, ok := table[method]; ok {
			once.Do(func() {
				if cfg.MaxRetry <= 0 {
					maxRetry = 1
				} else {
					maxRetry = cfg.MaxRetry
				}
				if cfg.RetryTimeout > 0 {
					retryTimeout = cfg.RetryTimeout
				}
				hystrix.ConfigureCommand(method, hystrix.CommandConfig{
					Timeout:               cfg.Timeout,
					MaxConcurrentRequests: cfg.MaxConcurrentRequests,

					SleepWindow: cfg.SleepWindow,

					RequestVolumeThreshold: cfg.RequestVolumeThreshold,

					ErrorPercentThreshold: cfg.ErrorPercentThreshold,
				})
			})
		} else {
			once.Do(func() {
				hystrix.ConfigureCommand(method, hystrix.CommandConfig{
					Timeout:                int(3 * time.Second),
					MaxConcurrentRequests:  3,
					SleepWindow:            5000,
					RequestVolumeThreshold: 20,
					ErrorPercentThreshold:  30,
				})
			})
		}
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
