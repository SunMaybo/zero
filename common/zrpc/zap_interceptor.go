package zrpc

import (
	"context"
	"github.com/SunMaybo/zero/common/zlog"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

func UnaryLoggerClientInterceptor() grpc.UnaryClientInterceptor {

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		startTime := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		logger := zlog.WithContext(ctx)
		if err != nil {
			logger.Errorw("grpc call error", zap.String("method", method), zap.String("status", "Failed"), zap.String("err", err.Error()), zap.Duration("cost", time.Since(startTime)))
		} else {
			logger.Infow("grpc call success", zap.String("method", method), zap.String("status", "Ok"), zap.Duration("cost", time.Since(startTime)))
		}
		return err
	}
}
func StreamLoggerClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		startTime := time.Now()
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		logger := zlog.WithContext(ctx)
		if err != nil {
			logger.Errorw("grpc call error", zap.String("method", method), zap.String("status", "Failed"), zap.String("err", err.Error()), zap.Duration("cost", time.Since(startTime)))
		} else {
			logger.Infow("grpc call success", zap.String("method", method), zap.String("status", "Ok"), zap.Duration("cost", time.Since(startTime)))
		}
		return clientStream, err
	}
}
func UnaryLoggerServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		result, err := handler(ctx, req)
		logger := zlog.WithContext(ctx)
		if err != nil {
			logger.Errorw("grpc call error", zap.String("method", info.FullMethod), zap.String("status", "Failed"), zap.String("err", err.Error()), zap.Duration("cost", time.Since(startTime)))
		} else {
			logger.Infow("grpc call success", zap.String("method", info.FullMethod), zap.String("status", "Ok"), zap.Duration("cost", time.Since(startTime)))
		}
		return result, err
	}
}
func StreamLoggerServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()
		err := handler(srv, stream)
		logger := zlog.WithContext(stream.Context())
		if err != nil {
			logger.Errorw("grpc call error", zap.String("method", info.FullMethod), zap.String("status", "Failed"), zap.String("err", err.Error()), zap.Duration("cost", time.Since(startTime)))
		} else {
			logger.Infow("grpc call success", zap.String("method", info.FullMethod), zap.String("status", "Ok"), zap.Duration("cost", time.Since(startTime)))
		}
		return err
	}
}
