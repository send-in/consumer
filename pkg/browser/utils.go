package browser

import (
	logger "consumer/pkg/log"
	"context"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/input"
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

func metaEnter() chromedp.Action {
    return chromedp.ActionFunc(func(ctx context.Context) error {
		if err := input.DispatchKeyEvent(input.KeyDown).
			WithKey("Meta").
			WithCode("MetaLeft").
			WithWindowsVirtualKeyCode(91).
			WithNativeVirtualKeyCode(91).
			Do(ctx); 
			err != nil {
			return err
		}

		if err := input.DispatchKeyEvent(input.KeyDown).
			WithKey("Enter").
			WithCode("Enter").
			WithWindowsVirtualKeyCode(13).
			WithNativeVirtualKeyCode(13).
			WithModifiers(input.ModifierMeta).
			Do(ctx); err != nil {
			return err
		}

		if err := input.DispatchKeyEvent(input.KeyUp).
			WithKey("Enter").
			WithCode("Enter").
			WithWindowsVirtualKeyCode(13).
			WithNativeVirtualKeyCode(13).
			WithModifiers(input.ModifierMeta).
			Do(ctx); err != nil {
			return err
		}

		return input.DispatchKeyEvent(input.KeyUp).
			WithKey("Meta").
			WithCode("MetaLeft").
			WithWindowsVirtualKeyCode(91).
			WithNativeVirtualKeyCode(91).
			Do(ctx)
	})
}

func setCSRFToken(token string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		return network.SetExtraHTTPHeaders(
			network.Headers{
				"csrf-token": strings.Trim(token, `"`),
			},
		).Do(ctx)
	})
}

func captureStatus(
	ctx context.Context, 
	url string, 
	timeout time.Duration,
) (bool, error) {
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
		case <-time.After(timeout):
			return false, context.DeadlineExceeded
	}
}