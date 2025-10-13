package browser

import (
	logger "consumer/pkg/log"
	"context"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func setCookie(name, value, domain, path string, httpOnly, secure bool) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
		err := network.SetCookie(name, value).
			WithExpires(&expr).
			WithDomain(domain).
			WithPath(path).
			WithHTTPOnly(httpOnly).
			WithSecure(secure).
			Do(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}

func setUserAgent(agent string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		err := emulation.SetUserAgentOverride(agent).Do(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}

func captureStatus(ctx context.Context, url string) (bool, error) {
	var requestId network.RequestID
	done := make(chan bool, 1)

	chromedp.ListenTarget(ctx, func(event any) {
		switch event := event.(type) {
			case *network.EventRequestWillBeSent:
				req := event.Request
				if strings.Contains(req.URL, url) {
					requestId = event.RequestID
				}

			case *network.EventResponseReceived:
				if event.RequestID == requestId && event.Response.Status == 200 {
					select {
						case done <- true:
							logger.Success("✅ Request successful! ID: %s", event.RequestID)
						default:
					}
				}
		}
	})

	select {
		case <-done:
			return true, nil
		case <-time.After(2*time.Second):
			return false, context.DeadlineExceeded			
	}
}

func (browser *Browser) SendTemp() (bool, error) {

	testURL := "https://reqres.in/"
	requestURL := "https://reqres.in/api/users?page=2"

	ctx, cancel := context.WithTimeout(browser.Context, 10*time.Second)
	defer cancel()

	err := chromedp.Run(
		ctx,
		chromedp.Navigate(testURL),
		chromedp.Evaluate(`
			(function() {
				let btn = document.createElement('button');
				btn.id = 'testBtn';
				btn.innerText = '`+"Message.message"+`';
				btn.onclick = function() {
					fetch('`+requestURL+`');
				};
				document.body.appendChild(btn);
			})()
		`, nil),

		chromedp.WaitVisible(`#testBtn`, chromedp.ByID),
		chromedp.Click(`#testBtn`, chromedp.ByID),
	)

	if err != nil {
		return false, err
	}

	status, err := captureStatus(
		browser.Context,
		requestURL,
	)
	if err != nil {
		return false, err
	}

	return status, nil
}
