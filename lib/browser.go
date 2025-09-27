package lib

import (
	"consumer/mq"
	"context"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func SetCookie(name, value, domain, path string, httpOnly, secure bool) chromedp.Action {
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

func (browser *Browser) Send(message mq.Message) error {
	fullMessage := message.Message + "\u2060"

	return chromedp.Run(browser.Context,
		SetCookie("li_at", message.Token, ".linkedin.com", "/", true, true),
		chromedp.Navigate(message.Reciever),

		chromedp.WaitVisible(`a.btn-primary.btn-sm.message-cta`, chromedp.ByQuery),
		chromedp.Click(`a.btn-primary.btn-sm.message-cta`, chromedp.ByQuery),

		chromedp.WaitVisible(`textarea#messaging-reply`, chromedp.ByQuery),
		chromedp.Click(`textarea#messaging-reply`, chromedp.ByQuery),
		chromedp.SendKeys(`textarea#messaging-reply`, fullMessage, chromedp.ByQuery),

		chromedp.Sleep(2*time.Second), // safety delay
		chromedp.WaitVisible(`button.message-send`, chromedp.ByQuery),
		chromedp.Click(`button.message-send`, chromedp.ByQuery),
	)
}
