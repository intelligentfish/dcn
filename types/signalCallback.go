package types

import "os"

// SignalCallback signal callback
type SignalCallback func(sig os.Signal)
