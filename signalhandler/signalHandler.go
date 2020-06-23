package signalhandler

import (
	"github.com/intelligentfish/dcn/types"
	"os"
)

// SignalHandler signal processor
type SignalHandler struct {
	callback types.SignalCallback
}

// New factory method
func New(callback types.SignalCallback) types.ISignalHandler {
	return &SignalHandler{callback: callback}
}

// OnSignal process signal
func (object *SignalHandler) OnSignal(sig os.Signal) {
	object.callback(sig)
}
