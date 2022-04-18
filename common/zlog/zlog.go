package zlog

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	LOGGER = NewLogger(false)
	zap.ReplaceGlobals(LOGGER)
	S = LOGGER.Sugar()
}

func InitLogger(production bool) {
	once.Do(func() {
		LOGGER = NewLogger(production)
		zap.ReplaceGlobals(LOGGER)
		S = LOGGER.Sugar()
	})
}
func NewLogger(production bool) *zap.Logger {
	level := zap.AtomicLevel{}
	if production {
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	} else {
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}
	cfg := zap.Config{
		Encoding: "console",
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Level:            level,
		Development:      !production,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			FunctionKey:    zapcore.OmitKey,
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	if logger, err := cfg.Build(); err != nil {
		panic(err)
	} else {
		return logger
	}
}
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := WithContext(c)
		t := time.Now()
		c.Next()
		latency := time.Since(t).String()
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		switch {
		case statusCode >= 400 && statusCode <= 499:
			{
				logger.Warn("[GIN]",
					zap.Int("statusCode", statusCode),
					zap.String("latency", latency),
					zap.String("clientIP", clientIP),
					zap.String("method", method),
					zap.String("path", path),
					zap.String("error", c.Errors.String()),
					zap.String("query", query),
				)
			}
		case statusCode >= 500:
			{
				logger.Error("[GIN]",
					zap.Int("statusCode", statusCode),
					zap.String("latency", latency),
					zap.String("clientIP", clientIP),
					zap.String("method", method),
					zap.String("path", path),
					zap.String("error", c.Errors.String()),
					zap.String("query", query),
				)
			}
		default:
			logger.Info("[GIN]",
				zap.Int("statusCode", statusCode),
				zap.String("latency", latency),
				zap.String("clientIP", clientIP),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("error", c.Errors.String()),
				zap.String("query", query),
			)
		}
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
				return S
			}

		default:
			return S
		}

	}
	return S
}
