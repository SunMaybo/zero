package zrpc

import (
	"context"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"reflect"
	"strings"
)

type ValidatorInterceptor struct {
	Interceptor grpc.UnaryServerInterceptor
}

type contextKey interface {
}

// NewValidatorInterceptor 参数检验rpc中间件
func NewValidatorInterceptor() *ValidatorInterceptor {
	validate := validator.New()
	_ = validate.RegisterValidationCtx("required_if_op", func(ctx context.Context, fl validator.FieldLevel) bool {
		if op, ok := ctx.Value("request_op").(string); ok {
			for _, v := range strings.Split(fl.Param(), " ") {
				if v == op {
					return hasValue(fl)
				}
			}
		}
		return true
	})
	_ = validate.RegisterValidationCtx("required_if_method", func(ctx context.Context, fl validator.FieldLevel) bool {
		if op, ok := ctx.Value("request_method").(string); ok {
			for _, v := range strings.Split(fl.Param(), " ") {
				if v == op {
					return hasValue(fl)
				}
			}
		}
		return true
	})
	// 注册手机号码校验规则，仅限国内手机号！！！国外手机号不要使用该tag！！！！！！
	_ = validate.RegisterValidation("mobile", isMobile)
	// 注册json数组校验规则
	_ = validate.RegisterValidation("strings", isStringSlice)
	_ = validate.RegisterValidation("numbers", isNumberSlice)
	_ = validate.RegisterValidation("chinese", isChineseChar)

	return &ValidatorInterceptor{
		Interceptor: RPCValidatorInterceptor(validate),
	}
}

func RPCValidatorInterceptor(validate *validator.Validate) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		validateCtx := parseMethod(info.FullMethod)
		if err := validate.StructCtx(validateCtx, req); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); !ok {
				for _, err := range err.(validator.ValidationErrors) {
					return nil, err
				}
			}
		}

		return handler(ctx, req)
	}
}

func parseMethod(fullMethod string) context.Context {
	ctx := context.Background()
	if list := strings.Split(fullMethod, "/"); len(list) > 1 {
		ctx = context.WithValue(ctx, contextKey("request_method"), list[len(list)-1])
	}

	method := strings.ToLower(fullMethod)
	methodList := []string{"create", "update", "delete", "save", "add", "get", "list", "rename"}

	for _, m := range methodList {
		if strings.Contains(method, m) {
			ctx = context.WithValue(ctx, contextKey("request_op"), m)
			return ctx
		}
	}
	return ctx
}

func hasValue(fl validator.FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		return !field.IsZero()
	}
}
