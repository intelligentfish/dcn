package main

import (
	"github.com/intelligentfish/dcn/app"
	"github.com/intelligentfish/dcn/log"
	"github.com/intelligentfish/dcn/srvGroup"
	"go.uber.org/zap"
	"os"
	"syscall"
)

func main() {
	log.Inst().Info("started", zap.String("app", app.Inst().Name()))
	var err error
	if err = srvGroup.Inst().StartAll(); nil != err {
		log.Inst().Error(err.Error())
		return
	}
	app.Inst().RegisterSignalCallback([]os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP},
		func(sig os.Signal) {
			srvGroup.Inst().StopAll()
		}).WaitSignal()
	log.Inst().Info("exited", zap.String("app", app.Inst().Name()))
}
