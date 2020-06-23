package srvGroup

import (
	"github.com/intelligentfish/dcn/app"
	"github.com/intelligentfish/dcn/define"
	"github.com/intelligentfish/dcn/serviceSorter"
	"github.com/intelligentfish/dcn/types"
	"os"
	"sync"
)

var (
	srvGroupInst *SrvGroup // service group instance
	srvGroupOnce sync.Once // service group once
)

// SrvGroup service group
type SrvGroup struct {
	sync.RWMutex
	services map[define.ServiceType]types.IService // service map
}

// new factory method
func new() *SrvGroup {
	return &SrvGroup{
		services: make(map[define.ServiceType]types.IService),
	}
}

// OnSignal signal processor
func (object *SrvGroup) OnSignal(sig os.Signal) {
	if !app.Inst().IsStopSignal(sig) {
		return
	}
	object.StopAll()
}

// AddSrv add service to group
func (object *SrvGroup) AddSrv(srv types.IService) *SrvGroup {
	object.Lock()
	defer object.Unlock()
	object.services[srv.GetServiceType()] = srv
	return object
}

// StartAll start all services
func (object *SrvGroup) StartAll() (err error) {
	services := make([]types.IService, 0, len(object.services))
	object.RLock()
	defer object.RUnlock()
	for _, srv := range object.services {
		if define.StartupPriorityUnknown == srv.GetStartupPriority() {
			continue
		}
		services = append(services, srv)
	}
	serviceSorter.New(services).SortByStartupPriority().Foreach(func(srv types.IService) (ok bool) {
		err = srv.Start()
		return nil == err
	})
	return
}

// StopAll stop all services
func (object *SrvGroup) StopAll() {
	services := make([]types.IService, 0, len(object.services))
	object.RLock()
	defer object.RUnlock()
	for _, srv := range object.services {
		if define.ShutdownPriorityUnknown == srv.GetShutdownPriority() {
			continue
		}
		services = append(services, srv)
	}
	serviceSorter.New(services).SortByShutdownPriority().Foreach(func(srv types.IService) (ok bool) {
		srv.Stop()
		return true
	})
}

// Inst singleton
func Inst() *SrvGroup {
	srvGroupOnce.Do(func() {
		srvGroupInst = new()
	})
	return srvGroupInst
}
