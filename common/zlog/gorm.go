package zlog

import (
	"context"
	"errors"
	"fmt"
	loggerGorm "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

func NewGormLogger(level loggerGorm.LogLevel) loggerGorm.Interface {
	config := loggerGorm.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  level,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	}
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = loggerGorm.Green + "%s\n" + loggerGorm.Reset + loggerGorm.Green + "[info] " + loggerGorm.Reset
		warnStr = loggerGorm.BlueBold + "%s\n" + loggerGorm.Reset + loggerGorm.Magenta + "[warn] " + loggerGorm.Reset
		errStr = loggerGorm.Magenta + "%s\n" + loggerGorm.Reset + loggerGorm.Red + "[error] " + loggerGorm.Reset
		traceStr = loggerGorm.Green + "%s\n" + loggerGorm.Reset + loggerGorm.Yellow + "[%.3fms] " + loggerGorm.BlueBold + "[rows:%v]" + loggerGorm.Reset + " %s"
		traceWarnStr = loggerGorm.Green + "%s " + loggerGorm.Yellow + "%s\n" + loggerGorm.Reset + loggerGorm.RedBold + "[%.3fms] " + loggerGorm.Yellow + "[rows:%v]" + loggerGorm.Magenta + " %s" + loggerGorm.Reset
		traceErrStr = loggerGorm.RedBold + "%s " + loggerGorm.MagentaBold + "%s\n" + loggerGorm.Reset + loggerGorm.Yellow + "[%.3fms] " + loggerGorm.BlueBold + "[rows:%v]" + loggerGorm.Reset + " %s"
	}

	return &logger{
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

type logger struct {
	loggerGorm.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode log mode
func (l *logger) LogMode(level loggerGorm.LogLevel) loggerGorm.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= loggerGorm.Info {
		WithContext(ctx).Infof(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= loggerGorm.Warn {
		WithContext(ctx).Warnf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= loggerGorm.Error {
		WithContext(ctx).Errorf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= loggerGorm.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= loggerGorm.Error && (!errors.Is(err, loggerGorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			WithContext(ctx).Errorf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			WithContext(ctx).Errorf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= loggerGorm.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			WithContext(ctx).Errorf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			WithContext(ctx).Errorf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == loggerGorm.Info:
		sql, rows := fc()
		if rows == -1 {
			WithContext(ctx).Errorf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			WithContext(ctx).Errorf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
