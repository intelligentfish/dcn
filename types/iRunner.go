package types

import (
	"context"
	"sync"
)

// IRunner runner interface
type IRunner interface {
	// Run method
	// ctx context
	// wg sync.WaitGroup
	// errCh error channel
	Run(ctx context.Context, wg *sync.WaitGroup, errCh chan<- error)
}
