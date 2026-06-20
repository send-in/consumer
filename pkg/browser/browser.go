package browser

import (
	mq "consumer/internal/queue"
	logger "consumer/pkg/log"

	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/chromedp/chromedp"
)

func NewBrowser(port int) *Browser {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.WindowSize(1366, 768),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		// chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.Flag("remote-debugging-port", strconv.Itoa(port)),
		chromedp.Flag("remote-debugging-address", "127.0.0.1"),
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	browserCtx, browserCancel := chromedp.NewContext(allocCtx)

	return &Browser{
		browserCtx,
		browserCancel,
		allocCancel,
		port,
	}
}

func (browser *Browser) Send(message mq.Message) (bool, error) {
	// fullMessage := message.Message + "\u2060"

	setupCtx, cancel := context.WithTimeout(browser.Context, 15*time.Second)
	defer cancel()

	err := chromedp.Run(
		setupCtx,
		setUserAgent(message.UserAgent),
		setCookie("li_at", message.Token, ".www.linkedin.com", "/", true, true),
		setCookie("JSESSIONID", message.JSession, ".www.linkedin.com", "/", true, true),
		chromedp.Navigate("https://www.linkedin.com/messaging/compose/"),
		chromedp.WaitReady("body"),
	)

	if err != nil {
		return false, err
	}

	statusCh := make(chan bool, 1)
	errCh := make(chan error, 1)

	go func() {
		status, err := captureStatus(
			browser.Context,
			SENDURL,
			30*time.Second,
		)

		if err != nil {
			errCh <- err
			return
		}

		statusCh <- status
	}()

	agentCtx, agentCancel := context.WithTimeout(
		context.Background(),
		45*time.Second,
	)
	defer agentCancel()

	cmd := exec.CommandContext(
		agentCtx,
		"agent-browser",
		"--cdp", strconv.Itoa(browser.Port),
		"batch",

		"wait --load networkidle",
		"keyboard type "+message.Receiver,
		"wait 1000",
		"press Enter",

		"snapshot -i",
		// "find role textbox click",
		// "keyboard type "+fullMessage,
		// "keypress Meta+Enter",
	)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	logger.Info("Executing: %s", cmd.String())

	if err := cmd.Run(); err != nil {
		logger.Error("STDOUT:\n%s", stdout.String())
		logger.Error("STDERR:\n%s", stderr.String())

		return false, fmt.Errorf(
			"agent-browser batch: %w\nstdout:\n%s\nstderr:\n%s",
			err,
			stdout.String(),
			stderr.String(),
		)
	}

	logger.Info("STDOUT:\n%s", stdout.String())
	logger.Info("STDERR:\n%s", stderr.String())

	select {
	case status := <-statusCh:
		return status, nil
	case err := <-errCh:
		return false, err
	case <-time.After(30 * time.Second):
		return false, context.DeadlineExceeded
	}
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

func (browser *Browser) Close() {
	browser.BrowserCancel()
	browser.AllocCancel()
}
