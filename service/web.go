package service

import (
	"context"
	"net"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/intelligentfish/dcn/define"
	"github.com/intelligentfish/dcn/log"
	"github.com/intelligentfish/dcn/serviceGroup"
	"github.com/intelligentfish/dcn/types"
)

var (
	webSrvInst *WebSrv
	webSrvOnce sync.Once
)

// WebSrv web service
type WebSrv struct {
	*BaseSrv
	lAddr                      string
	engine                     *gin.Engine
	healthCheckHandlerSet      map[types.IHealthCheck]struct{}
	healthCheckHandlerSetMutex sync.RWMutex
	ChildHandler               types.IWebHandler
}

// NewWebSrv factory method
func NewWebSrv(options ...WebSrvOption) *WebSrv {
	object := &WebSrv{
		BaseSrv: NewSrvBase("WebSrv",
			define.ServiceTypeWeb,
			define.StartupPriorityWeb,
			define.ShutdownPriorityWeb),
		engine:                gin.New(),
		healthCheckHandlerSet: make(map[types.IHealthCheck]struct{}, 0),
	}
	object.ChildRunner = object
	object.Use(options...)
	object.engine.Use(gin.LoggerWithWriter(object), gin.Recovery())
	if nil != object.ChildHandler {
		object.ChildHandler.UseEngine(object.engine)
	} else {
		object.UseEngine(object.engine)
	}
	return object
}

// Write sink gin log to zap
func (object *WebSrv) Write(p []byte) (n int, err error) {
	log.Inst().Info(string(p))
	n = len(p)
	return
}

// AddHealthChecker add health checker
func (object *WebSrv) AddHealthChecker(handler types.IHealthCheck) *WebSrv {
	object.healthCheckHandlerSetMutex.Lock()
	defer object.healthCheckHandlerSetMutex.Unlock()
	object.healthCheckHandlerSet[handler] = struct{}{}
	return object
}

// UseEngine use web engine
func (object *WebSrv) UseEngine(engine *gin.Engine) {
	engine.GET("/api/v1/health", func(ctx *gin.Context) {
		object.healthCheckHandlerSetMutex.RLock()
		defer object.healthCheckHandlerSetMutex.RUnlock()
		var err error
		for handler := range object.healthCheckHandlerSet {
			if err = handler.OnHealthCheck(); nil != err {
				log.Inst().Error(err.Error())
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
		}
		ctx.Status(http.StatusOK)
	})
}

// Use use options for web service
func (object *WebSrv) Use(options ...WebSrvOption) *WebSrv {
	for _, option := range options {
		option(object)
	}
	return object
}

// Run run method
func (object *WebSrv) Run(ctx context.Context, wg *sync.WaitGroup, errCh chan<- error) {
	ln, err := net.Listen("tcp", object.lAddr)
	if nil != err {
		errCh <- err
		return
	}
	go func() {
		<-ctx.Done()
		if err = ln.Close(); nil != err && err != http.ErrServerClosed {
			log.Inst().Info(err.Error())
		}
	}()
	if err = object.engine.RunListener(ln); nil != err {
		log.Inst().Info(err.Error())
	}
}

// WebSrvInst web service singleton
func WebSrvInst() *WebSrv {
	webSrvOnce.Do(func() {
		//TODO config read listen addr
		opt := WebSrvListenAddrOption(":80")
		webSrvInst = NewWebSrv(opt)
	})
	return webSrvInst
}

// init init method
func init() {
	serviceGroup.Inst().AddSrv(WebSrvInst())
}
