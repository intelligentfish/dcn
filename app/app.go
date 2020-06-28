package app

import (
	"github.com/intelligentfish/dcn/signalHandler"
	"github.com/intelligentfish/dcn/types"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"sync/atomic"
	"syscall"
)

var (
	appInst *App      // app instance
	appOnce sync.Once // app once
)

// App application
type App struct {
	stopped            int32
	name               string
	stopSigSet         map[os.Signal]struct{}
	stopSigSetMutex    sync.RWMutex
	sigHandlerMap      map[os.Signal][]types.ISignalHandler
	sigHandlerMapMutex sync.RWMutex
}

// new factory method
func new() *App {
	object := &App{
		sigHandlerMap: make(map[os.Signal][]types.ISignalHandler, 0),
		stopSigSet:    make(map[os.Signal]struct{}, 0),
	}
	object.name = filepath.Base(os.Args[0])
	object.stopSigSet[syscall.SIGINT] = struct{}{}
	object.stopSigSet[syscall.SIGTERM] = struct{}{}
	object.stopSigSet[syscall.SIGHUP] = struct{}{}
	return object
}

// notifySignal notify signal arrival
func (object *App) notifySignal(sig os.Signal) {
	object.sigHandlerMapMutex.RLock()
	defer object.sigHandlerMapMutex.RUnlock()
	if arr, ok := object.sigHandlerMap[sig]; ok {
		for _, handler := range arr {
			handler.OnSignal(sig)
		}
	}
}

// Name app name
func (object *App) Name() string {
	return object.name
}

// IsStopped return true if app need stop
func (object *App) IsStopped() bool {
	return 1 == atomic.LoadInt32(&object.stopped)
}

// AddStopSignal add stop signal
func (object *App) AddStopSignal(sig os.Signal) *App {
	object.stopSigSetMutex.Lock()
	object.stopSigSetMutex.Unlock()
	object.stopSigSet[sig] = struct{}{}
	return object
}

// RegisterSignalHandler register signal handler
func (object *App) RegisterSignalHandler(sigList []os.Signal,
	handler types.ISignalHandler) *App {
	object.sigHandlerMapMutex.Lock()
	defer object.sigHandlerMapMutex.Unlock()
	for _, sig := range sigList {
		if _, ok := object.sigHandlerMap[sig]; !ok {
			object.sigHandlerMap[sig] = []types.ISignalHandler{handler}
		} else {
			object.sigHandlerMap[sig] = append(object.sigHandlerMap[sig], handler)
		}
	}
	return object
}

// RegisterSignalCallback register signal callback
func (object *App) RegisterSignalCallback(sigList []os.Signal,
	callback types.SignalCallback) *App {
	object.RegisterSignalHandler(sigList, signalHandler.New(callback))
	return object
}

// IsStopSignal whether signal in stop signal set
func (object *App) IsStopSignal(sig os.Signal) (ok bool) {
	object.stopSigSetMutex.RLock()
	defer object.stopSigSetMutex.RUnlock()
	_, ok = object.stopSigSet[sig]
	return
}

// WaitSignal wait for signal in list
func (object *App) WaitSignal(sigList ...os.Signal) {
	sigCh := make(chan os.Signal, 32)
	signal.Notify(sigCh, sigList...)
loop:
	for sig := range sigCh {
		object.notifySignal(sig)
		if object.IsStopSignal(sig) {
			break loop
		}
	}
	atomic.StoreInt32(&object.stopped, 1)
	signal.Stop(sigCh)
	close(sigCh)
}

// Inst singleton
func Inst() *App {
	appOnce.Do(func() {
		appInst = new()
	})
	return appInst
}
