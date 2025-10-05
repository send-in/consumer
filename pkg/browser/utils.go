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

	err := network.Enable().Do(ctx)
	if  err != nil {
		return false, err
	}

	chromedp.ListenTarget(ctx, func(event interface{}) {
		switch event := event.(type) {
			case *network.EventRequestWillBeSent:
				req := event.Request
				if strings.Contains(req.URL, url) {
					logger.Info("Request ID: %s", event.RequestID)
					requestId = event.RequestID
				}

			case *network.EventResponseReceived:
				if event.RequestID == requestId && event.Response.Status == 200 {
					select{
						case done<-true:
							logger.Success("✅ Request successful! ID: %s", event.RequestID)
						default:
					}
				}
		}
	})

	select{
		case <-done:
			return true, nil
		case <-time.After(10*time.Second):
			return false, context.DeadlineExceeded			
	}
}
