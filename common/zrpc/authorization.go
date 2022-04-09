package zrpc

import (
	"context"
	"errors"
	"github.com/SunMaybo/zero/common/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

type TokenAuth struct {
	token     string
	enableTLS bool
}

func NewTokenAuth(token string, enableTLS bool) *TokenAuth {
	return &TokenAuth{
		token:     token,
		enableTLS: enableTLS,
	}
}
func (t *TokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + t.token,
	}, nil
}
func (t *TokenAuth) RequireTransportSecurity() bool {
	return t.enableTLS
}

func UnaryJWTServerInterceptor(secretKey string, filterMethods []RpcMethod) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		for _, method := range filterMethods {
			if strings.HasSuffix(info.FullMethod, string(method)) {
				return handler(ctx, req)
			}
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("no metadata")
		}
		if md["authorization"] == nil || len(md["authorization"]) <= 0 {
			return nil, errors.New("no metadata")
		}
		if payload, err := jwt.NewParseToken(secretKey).Verify(md["authorization"][0]); err != nil {
			return nil, err
		} else {
			ctx = context.WithValue(ctx, "payload", payload)
			return handler(ctx, req)
		}
	}
}

//func UnaryJWTStreamServerInterceptor(secretKey string, filterMethods []string) grpc.StreamServerInterceptor {
//	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
//		for _, method := range filterMethods {
//			if method == info.FullMethod {
//				return handler(srv, ss)
//			}
//		}
//		md, ok := metadata.FromIncomingContext(ss.Context())
//		if !ok {
//			return errors.New("no metadata")
//		}
//		if md["authorization"] == nil || len(md["authorization"]) <= 0 {
//			return errors.New("no metadata")
//		}
//		if payload, err := jwt.NewParseToken(secretKey).Verify(md["authorization"][0]); err != nil {
//			return err
//		} else {
//
//			ss. = context.WithValue(ss.Context(), "payload", payload)
//			return handler(ctx, req)
//		}
//	}
//}
