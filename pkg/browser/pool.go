package browser

import (
	"sync"
	"time"
)

func CreatePool(size int) *BrowserPool {
	pool := make(chan *Browser, size)
	for range size {
		pool <- NewBrowser()
	}

	return &BrowserPool{
		Pool:   pool,
		Size:   size,
		closed: false,
		mu:     &sync.RWMutex{},
	}
}

func (browsers *BrowserPool) Acquire() *Browser {
	browsers.mu.RLock()
	defer browsers.mu.RUnlock()

	if browsers.closed {
		return nil
	}

	select {
	case browser := <-browsers.Pool:
		return browser
	case <-time.After(30 * time.Second): // Timeout after 30 seconds
		return nil
	}
}

func (browsers *BrowserPool) Release(browser *Browser) {
	browsers.mu.RLock()
	defer browsers.mu.RUnlock()

	if browsers.closed || browser == nil {
		if browser != nil {
			browser.Close()
		}
		return
	}

	select {
	case browsers.Pool <- browser:
		// Successfully returned to pool
	default:
		// Pool is full, close the browser
		browser.Close()
	}
}

// Close gracefully shuts down all browsers in the pool
// func (browsers *BrowserPool) Close() {
// 	browsers.mu.Lock()
// 	defer browsers.mu.Unlock()

// 	if browsers.closed {
// 		return
// 	}

// 	browsers.closed = true
// 	close(browsers.Pool)

// 	// Close all remaining browsers in the pool
// 	for browser := range browsers.Pool {
// 		browser.Close()
// 	}
// }

// Size returns the configured size of the pool
// func (browsers *BrowserPool) GetSize() int {
// 	return browsers.Size
// }

// Available returns the number of browsers currently available in the pool
// func (browsers *BrowserPool) Available() int {
// 	browsers.mu.RLock()
// 	defer browsers.mu.RUnlock()

// 	if browsers.closed {
// 		return 0
// 	}

// 	return len(browsers.Pool)
// }
