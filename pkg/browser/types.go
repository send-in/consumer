package browser

import (
	"context"
)

const(
	SENDURL = "https://www.linkedin.com/voyager/api/voyagerMessagingDashMessengerMessages?action=createMessage"
)

type Browser struct {
	Context       context.Context
	BrowserCancel context.CancelFunc
	AllocCancel   context.CancelFunc
}

type BrowserPool struct {
	Size   int
	Pool   chan *Browser
	Active bool
}