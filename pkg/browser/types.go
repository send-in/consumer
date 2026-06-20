package browser

import (
	"context"
	"sync/atomic"
)

const(
	SENDURL = "https://www.linkedin.com/voyager/api/voyagerMessagingDashMessengerMessages?action=createMessage"
)

type Browser struct {
	Context       context.Context
	BrowserCancel context.CancelFunc
	AllocCancel   context.CancelFunc
	Port		  int
}

type BrowserPool struct {
	Size   int
	Pool   chan *Browser
	Dead   chan *Browser
	Active *atomic.Bool
}