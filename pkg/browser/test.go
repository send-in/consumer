package browser

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

func (browser *Browser) SendTemp() (bool, error) {

	testURL := "https://reqres.in/"
	requestURL := "https://reqres.in/api/users?page=2"

	ctx, cancel := context.WithTimeout(browser.Context, 10*time.Second)
	defer cancel()

	err := chromedp.Run(
		ctx,
		// chromedp.Sleep(11*time.Second),
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
