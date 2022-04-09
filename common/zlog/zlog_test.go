package zlog

import (
	"go.uber.org/zap"
	"testing"
)

func TestZlog(t *testing.T) {
	InitLogger(false)
	LOGGER.Info("当前执行消息", zap.String("name", "maybo"))
}
