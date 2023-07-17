package zlog

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

var LOGGER *zap.Logger
var S *zap.SugaredLogger

var once = sync.Once{}

func init() {
	LOGGER = NewLogger(false, "")
	zap.ReplaceGlobals(LOGGER)
	S = LOGGER.Sugar()
}

func InitLogger(production bool, fileName string) {
	once.Do(func() {
		LOGGER = NewLogger(production, fileName)
		zap.ReplaceGlobals(LOGGER)
		S = LOGGER.Sugar()
	})
}
func NewLogger(production bool, fileName string) *zap.Logger {
	level := zap.AtomicLevel{}
	if production {
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	} else {
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}
	//cfg := zap.Config{
	//	Encoding:    "console",
	//	OutputPaths: []string{"stderr"},
	//	EncoderConfig: zapcore.EncoderConfig{
	//		MessageKey:  "message",
	//		TimeKey:     "time",
	//		EncodeTime:  zapcore.RFC3339TimeEncoder,
	//		LevelKey:    "level",
	//		EncodeLevel: zapcore.CapitalColorLevelEncoder,
	//	},
	//	Level: level,
	//}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	var tee zapcore.Core
	if fileName != "" {
		fileSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28,    //days
			Compress:   false, // disabled by default
		})
		consoleSyncer := zapcore.AddSync(os.Stdout)
		tee = zapcore.NewTee(zapcore.NewCore(encoder, consoleSyncer, level), zapcore.NewCore(encoder, fileSyncer, level))
	} else {
		consoleSyncer := zapcore.AddSync(os.Stdout)
		tee = zapcore.NewTee(zapcore.NewCore(encoder, consoleSyncer, level))
	}
	logger := zap.New(tee, zap.AddCaller())
	return logger
}
func GinLogger(c context.Context, startTime time.Time, ginCtx *gin.Context) {
	logger := WithContext(c)
	latency := time.Since(startTime).String()
	method := ginCtx.Request.Method
	statusCode := ginCtx.Writer.Status()
	path := ginCtx.Request.URL.Path
	switch {
	case statusCode >= 400 && statusCode <= 499:
		{
			logger.Warnf("%s %s %d %s", method, path, statusCode, latency)
			if ginCtx.Errors != nil && ginCtx.Errors.ByType(gin.ErrorTypePrivate).String() != "" {
				logger.Warnf("request error: %s", ginCtx.Errors.ByType(gin.ErrorTypePrivate).String())
			}

		}
	case statusCode >= 500:
		{
			logger.Errorf("%s %s %d %s", method, path, statusCode, latency)
			if ginCtx.Errors != nil && ginCtx.Errors.ByType(gin.ErrorTypePrivate).String() != "" {
				logger.Errorf("request error: %s", ginCtx.Errors.ByType(gin.ErrorTypePrivate).String())
			}
		}
	default:
		logger.Infof("%s %s %d %s", method, path, statusCode, latency)
	}

}
func RecoveryWithZap() gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			logger := WithContext(c)
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.Abort()
					return
				}
				logger.Error("[Recovery from panic]",
					zap.Time("time", time.Now()),
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					zap.String("stack", string(debug.Stack())),
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
func WithContext(ctx context.Context) *zap.SugaredLogger {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		switch span.Context().(type) {
		case jaeger.SpanContext:
			span := span.Context().(jaeger.SpanContext)
			return S.With(zap.String("parent_id", span.ParentID().String()), zap.String("span_id", span.SpanID().String()), zap.String("trace_id", span.TraceID().String()))
		case zipkinot.SpanContext:
			span := span.Context().(zipkinot.SpanContext)
			if span.ParentID != nil {
				return S.With(zap.String("parent_id", span.ParentID.String()), zap.String("span_id", span.ID.String()), zap.String("trace_id", span.TraceID.String()))
			} else {
				return S.With(zap.String("span_id", span.ID.String()), zap.String("trace_id", span.TraceID.String()))
			}

		default:
			return S
		}

	}
	return S
}
