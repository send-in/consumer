package lib

import (
	"consumer/mq"
	"context"

	"github.com/chromedp/chromedp"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

type Browser struct {
	Context context.Context
	Cancel context.CancelFunc
}

func CreateBrowser(message *mq.Message) *Browser {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.UserAgent(message.UserAgent),
		chromedp.WindowSize(1920, 1080),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.DisableGPU,
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	browserCtx, cancel := chromedp.NewContext(allocCtx)
	
	return &Browser{
		browserCtx,
		cancel,
	}
}