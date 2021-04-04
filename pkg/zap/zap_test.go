package zap

import (
	"github.com/go-kratos/kratos/v2/log"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger, err := NewLogger(AddCaller())
	if err != nil {
		t.Error(err)
	}
	log.Debug(logger).Print("log", "test")
	log.Info(logger).Print("log", "test")
	log.Warn(logger).Print("log", "test")
	log.Error(logger).Print("log", "test")

}
