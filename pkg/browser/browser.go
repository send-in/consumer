package browser

import (
	mq "consumer/internal/queue"

	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func NewBrowser() *Browser {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.WindowSize(700, 700),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("ignore-certificate-errors", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	browserCtx, browserCancel := chromedp.NewContext(allocCtx)

	return &Browser{
		browserCtx,
		browserCancel,
		allocCancel,
	}
}

func (browser *Browser) Send(message mq.Message)  (bool, error) {
    ctx, cancel := context.WithTimeout(
        browser.Context,
        60*time.Second,
    )

    defer cancel()

	textbox := `//*[@contenteditable="true" and @role="textbox"]`
    composeURL := fmt.Sprintf(
        "https://www.linkedin.com/messaging/compose/?profileUrn=%s&recipient=%s&screenContext=NON_SELF_PROFILE_VIEW&interop=msgOverlay",
        url.QueryEscape(message.ProfileURN),
        url.QueryEscape(message.Recipient),
    )

	
    err := chromedp.Run(
		ctx,
		network.Enable(),
        setCookie("JSESSIONID", message.JSession, ".linkedin.com", "/", true, true),
        setCookie("li_at", message.Token, ".linkedin.com", "/", true, true),
        setCSRFToken(message.JSession),
        chromedp.Navigate(composeURL),
        chromedp.WaitVisible(textbox, chromedp.BySearch),
        chromedp.Click(textbox, chromedp.BySearch),
        chromedp.SendKeys(
			textbox, 
			message.Message, 
			chromedp.BySearch,
		),

		metaEnter(),
    )

	if err != nil {
		return false, err
	}

	status, err := captureStatus(
		browser.Context,
		SENDURL,
		60*time.Second,
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
		chromedp.Navigate("https://www.linkedin.com/messaging/"),
	)
}

func (browser *Browser) Close() {
	browser.BrowserCancel()
	browser.AllocCancel()
}
