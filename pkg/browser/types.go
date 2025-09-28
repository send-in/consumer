package browser

import (
	"context"
	"sync"
)

type Browser struct {
	Context       context.Context
	BrowserCancel context.CancelFunc
	AllocCancel   context.CancelFunc
}

type BrowserPool struct {
	Pool   chan *Browser
	Size   int
	closed bool
	mu     *sync.RWMutex
}
