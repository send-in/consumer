package browser

import (
	logger "consumer/pkg/log"
	"errors"
	"time"
)

func CreatePool(size int) (*BrowserPool, error) {
	pool := make(chan *Browser, size)
	for range size {
		pool <- NewBrowser()
	}
	
	browsers := &BrowserPool{
		size,
		pool,
		true,
	}
	
	if len(browsers.Pool) != browsers.Size {
		return nil, errors.New("browsers crashed") 
	}

	return browsers, nil
}

func (browsers *BrowserPool) Acquire() *Browser {
	select {
	case browser := <-browsers.Pool:
		return browser
	case <-time.After(30 * time.Second):
		return nil
	}
}

func (browsers *BrowserPool) Release(browser *Browser) {
	if browser == nil || !browser.IsAlive() {
        browser.Close()
        select {
		case browsers.Pool <- NewBrowser():
			logger.Warning("💀 Replaced crashed browser during Release()")
		default:
        }
    }else {
		select {
		case browsers.Pool <- browser:
		default:
			browser.Close()
		}
	}
}

func (browsers *BrowserPool) Close() {
	browsers.Active = false
	close(browsers.Pool)
	
	for browser := range browsers.Pool {
		browser.Close()
	}
}