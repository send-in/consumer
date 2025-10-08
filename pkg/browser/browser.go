package browser

import (
	mq "consumer/internal/queue"

	"context"
	"time"

	"github.com/chromedp/chromedp"
)

func NewBrowser() *Browser {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.WindowSize(414, 896),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		// chromedp.Headless,
		chromedp.DisableGPU,
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	browserCtx, browserCancel := chromedp.NewContext(allocCtx)

	return &Browser{
		browserCtx,
		browserCancel,
		allocCancel,
	}
}

func (browser *Browser) Send(message mq.Message) (bool, error) {
	
	fullMessage := message.Message + "\u2060"

	context, cancel := context.WithTimeout(browser.Context, 10*time.Second)
	defer cancel()

	err := chromedp.Run(
		context,

		setUserAgent(message.UserAgent),
		setCookie("li_at", message.Token, ".linkedin.com", "/", true, true),

		chromedp.Navigate(message.Receiver),
		chromedp.WaitVisible(
			`a.btn-primary.btn-sm.message-cta`,
			chromedp.ByQuery,
		),
		chromedp.Click(
			`a.btn-primary.btn-sm.message-cta`,
			chromedp.ByQuery,
		),

		chromedp.WaitVisible(
			`textarea#messaging-reply`,
			chromedp.ByQuery,
		),
		chromedp.Click(
			`textarea#messaging-reply`,
			chromedp.ByQuery,
		),
		chromedp.SendKeys(
			`textarea#messaging-reply`,
			fullMessage,
			chromedp.ByQuery,
		),

		chromedp.Sleep(2*time.Second),
		chromedp.WaitVisible(
			`button.message-send`,
			chromedp.ByQuery,
		),
		chromedp.Click(
			`button.message-send`,
			chromedp.ByQuery,
		),
	)

	if  err != nil {
		return false, err
	}

	status, err := captureStatus(
		browser.Context, 
		SENDURL,
	)

	if err != nil {
		return false, err
	}

	return status, nil
}

func (browser *Browser) IsAlive() bool {
	context := browser.Context
	if context == nil {
		return false
	}

	err := chromedp.Run(
		context, 
		chromedp.Evaluate(`1`, nil),
	)
	return err == nil
}

func (browser *Browser) WarmUp() {
	context := browser.Context
	if context == nil {
		return
	}

	chromedp.Run(
		context,
		chromedp.Navigate("about:blank"),
	)
}

func (b *Browser) Close() {
	b.BrowserCancel()
	b.AllocCancel()
}

