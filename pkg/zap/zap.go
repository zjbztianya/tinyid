package zap

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_        log.Logger = (*Logger)(nil)
	levelMap            = map[log.Level]zapcore.Level{
		log.LevelDebug: zapcore.DebugLevel,
		log.LevelInfo:  zapcore.InfoLevel,
		log.LevelWarn:  zapcore.WarnLevel,
		log.LevelError: zapcore.ErrorLevel,
	}
)

type Logger struct {
	log *zap.Logger
}

func AddCaller() zap.Option {
	return zap.AddCaller()
}

func AddCallerSkip(skip int) zap.Option {
	return zap.AddCallerSkip(skip)
}

func AddStacktrace(lvl log.Level) zap.Option {
	return zap.AddStacktrace(levelMap[lvl])
}

func NewLogger(opts ...zap.Option) (*Logger, error) {
	l, err := zap.NewProduction(opts...)
	if err != nil {
		return nil, err
	}
	return &Logger{log: l}, nil
}

func (l *Logger) Print(pairs ...interface{}) {
	level := pairs[1].(log.Level)
	zapLevel, ok := levelMap[level]
	if !ok {
		return
	}
	pairs = pairs[2:]
	msg := fmt.Sprint(pairs...)
	if ce := l.log.Check(zapLevel, msg); ce != nil {
		ce.Write()
	}
}
