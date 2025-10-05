package browser

import (
	logger "consumer/pkg/log"
	"time"
)

func CreatePool(size int) *BrowserPool {
	pool := make(chan *Browser, size)
	for range size {
		pool <- NewBrowser()
	}
	
	browsers := &BrowserPool{
		size,
		pool,
		true,
	}
	
	go browsers.Monitor()
	return browsers
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
	select {
	case browsers.Pool <- browser:
	default:
		browser.Close()
	}
}

func (browsers *BrowserPool) Close() {
	browsers.Active = false
	close(browsers.Pool)
	
	for browser := range browsers.Pool {
		browser.Close()
	}
}

func (browsers *BrowserPool) Monitor(){
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	
	// TODO: cleanup logic

	for browsers.Active {
		<-ticker.C
		
		required := browsers.Size
		actual := len(browsers.Pool)
		
		for range (required-actual){
			select {
			case browsers.Pool <- NewBrowser():
				logger.Info("♻️ Browser added to pool for self-healing")
			default:
			}
		}
	}
}