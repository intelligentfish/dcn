package main

import (
	"flag"
	"github.com/intelligentfish/dcn/agent/module/http_module"
	"github.com/intelligentfish/dcn/app"
	"github.com/intelligentfish/dcn/config"
	"github.com/intelligentfish/dcn/log"
	_ "github.com/intelligentfish/dcn/service"
	"github.com/intelligentfish/dcn/serviceGroup"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"os"
	"syscall"
)

// asAgent run as agent
func asAgent() (err error) {
	// check config, use task pool, concurrency run task
	L := lua.NewState()
	defer L.Close()
	http_module.Inject(L)
	for {
		// pull task
		// schedule task
		// execute task
		// loop
		L.SetGlobal("token", lua.LString("task-id"))
		err = L.DoString(``)
		if nil != err {
			log.Inst().Error(err.Error())
		}
	}
	return
}

// asServer run as server
func asServer() (err error) {
	if err = serviceGroup.Inst().StartAll(); nil != err {
		log.Inst().Error(err.Error())
		return
	}
	app.Inst().RegisterSignalCallback([]os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP},
		func(sig os.Signal) {
			serviceGroup.Inst().StopAll()
		})
	return
}

// main entry point
func main() {
	flag.BoolVar(&config.Inst().Agent, "agent", false, "run as agent")
	flag.Parse()
	log.Inst().Info("started", zap.String("app", app.Inst().Name()))
	var err error
	if config.Inst().Agent {
		err = asAgent()
	} else {
		err = asServer()
	}
	if nil != err {
		log.Inst().Error(err.Error())
		return
	}

	app.Inst().WaitSignal()
	log.Inst().Info("exited", zap.String("app", app.Inst().Name()))
}
