package log

import (
	"fmt"
	"github.com/intelligentfish/dcn/app"
	"go.uber.org/zap"
	"os"
	"sync"
)

var (
	logInst *Log      // log instance
	logOnce sync.Once // log once
)

// Log log instance
type Log struct {
	logger *zap.Logger
}

// NewLog factory method
func NewLog() *Log {
	var err error
	var logger *zap.Logger
	zapCfg := zap.NewProductionConfig()
	//TODO read from config
	zapCfg.Development = true
	zapCfg.DisableCaller = false
	zapCfg.OutputPaths = []string{
		"stdout",
		os.ExpandEnv(fmt.Sprintf("/var/log/%s.stdout.log", app.Inst().Name())),
	}
	zapCfg.ErrorOutputPaths = []string{
		"stderr",
		os.ExpandEnv(fmt.Sprintf("/var/log/%s.stderr.log", app.Inst().Name())),
	}
	if logger, err = zapCfg.Build(); nil != err {
		panic(err)
	}
	logger.Named(app.Inst().Name())
	return &Log{logger: logger}
}

// Inst singleton
func Inst() *zap.Logger {
	logOnce.Do(func() {
		logInst = NewLog()
	})
	return logInst.logger
}
