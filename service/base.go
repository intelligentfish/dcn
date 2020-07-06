package service

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/intelligentfish/dcn/define"
	"github.com/intelligentfish/dcn/types"
)

// BaseSrv base service
type BaseSrv struct {
	name             string
	srvType          define.ServiceType
	startupPriority  define.StartupPriority
	shutdownPriority define.ShutdownPriority
	startFlag        int32
	wg               sync.WaitGroup
	ctx              context.Context
	cancel           context.CancelFunc
	ChildRunner      types.IRunner
}

// NewSrvBase factory method
func NewSrvBase(name string,
	srvType define.ServiceType,
	startupPriority define.StartupPriority,
	shutdownPriority define.ShutdownPriority,
) *BaseSrv {
	object := &BaseSrv{
		name:             name,
		srvType:          srvType,
		startupPriority:  startupPriority,
		shutdownPriority: shutdownPriority,
	}
	object.ctx, object.cancel = context.WithCancel(context.Background())
	return object
}

// SetKeyValue set key value pair into context
func (object *BaseSrv) SetKeyValue(key, value interface{}) {
	object.ctx = context.WithValue(object.ctx, key, value)
}

// GetValue get value from context
func (object *BaseSrv) GetValue(key interface{}) (value interface{}) {
	return object.ctx.Value(key)
}

// GetServiceType get the service type
func (object *BaseSrv) GetServiceType() define.ServiceType {
	return object.srvType
}

// GetStartupPriority get the startup priority
func (object *BaseSrv) GetStartupPriority() define.StartupPriority {
	return object.startupPriority
}

// GetShutdownPriority get the shutdown priority
func (object *BaseSrv) GetShutdownPriority() define.ShutdownPriority {
	return object.shutdownPriority
}

// Name service name
func (object *BaseSrv) Name() string {
	return object.name
}

// Start start the service
func (object *BaseSrv) Start() (err error) {
	if !atomic.CompareAndSwapInt32(&object.startFlag, 0, 1) {
		return
	}
	object.wg.Add(1)
	errCh := make(chan error, 1)
	go func() {
		defer object.wg.Done()
		if nil == object.ChildRunner {
			object.Run(object.ctx, &object.wg, errCh)
		} else {
			object.ChildRunner.Run(object.ctx, &object.wg, errCh)
		}
	}()
	err = <-errCh
	return
}

// Stop stop the service
func (object *BaseSrv) Stop() {
	if !atomic.CompareAndSwapInt32(&object.startFlag, 1, 0) {
		return
	}
	object.cancel()
	object.wg.Wait()
}

// Run real service method
func (object *BaseSrv) Run(ctx context.Context,
	wg *sync.WaitGroup,
	errCh chan<- error) {
	errCh <- nil
}
