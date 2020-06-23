package types

import "os"

// ISignalHandler signal handler
type ISignalHandler interface {
	OnSignal(sig os.Signal) // signal processor
}
