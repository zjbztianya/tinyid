package zap

import (
	"bytes"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
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
	log  *zap.Logger
	pool *sync.Pool
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
	return &Logger{
		log: l,
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}, nil
}

func (l *Logger) Print(pairs ...interface{}) {
	level := pairs[1].(log.Level)
	zapLevel, ok := levelMap[level]
	if !ok {
		return
	}
	if len(pairs)%2 != 0 {
		pairs = append(pairs, "")
	}
	pairs = pairs[2:]
	buf := l.pool.Get().(*bytes.Buffer)
	for i := 0; i < len(pairs); i += 2 {
		fmt.Fprintf(buf, "%s=%v ", pairs[i], log.Value(pairs[i+1]))
	}
	if ce := l.log.Check(zapLevel, buf.String()); ce != nil {
		ce.Write()
	}
	buf.Reset()
	l.pool.Put(buf)
}
