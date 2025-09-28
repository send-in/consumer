package browser

import (
	"sync"
	"time"
)

// BrowserPoolWithTimeout extends BrowserPool with idle timeout functionality
type BrowserPoolWithTimeout struct {
	*BrowserPool
	idleTimeout time.Duration
	lastUsed    map[*Browser]time.Time
	lastUsedMu  sync.RWMutex
	stopCleanup chan bool
	wg          sync.WaitGroup
}

// CreatePoolWithTimeout creates a browser pool with idle timeout management
func CreatePoolWithTimeout(size int, idleTimeout time.Duration) *BrowserPoolWithTimeout {
	basePool := CreatePool(size)

	pool := &BrowserPoolWithTimeout{
		BrowserPool: basePool,
		idleTimeout: idleTimeout,
		lastUsed:    make(map[*Browser]time.Time),
		stopCleanup: make(chan bool),
	}

	// Start the cleanup goroutine
	pool.wg.Add(1)
	go pool.cleanupIdleBrowsers()

	return pool
}

// Acquire gets a browser from the pool and tracks its usage time
func (p *BrowserPoolWithTimeout) Acquire() *Browser {
	browser := p.BrowserPool.Acquire()
	if browser != nil {
		p.lastUsedMu.Lock()
		p.lastUsed[browser] = time.Now()
		p.lastUsedMu.Unlock()
	}
	return browser
}

// Release returns a browser to the pool and updates its last used time
func (p *BrowserPoolWithTimeout) Release(browser *Browser) {
	if browser != nil {
		p.lastUsedMu.Lock()
		p.lastUsed[browser] = time.Now()
		p.lastUsedMu.Unlock()
	}
	p.BrowserPool.Release(browser)
}

// Close shuts down the pool and cleanup goroutine
func (p *BrowserPoolWithTimeout) Close() {
	// Stop the cleanup goroutine
	close(p.stopCleanup)
	p.wg.Wait()

	// Close the base pool
	p.BrowserPool.Close()

	// Clean up tracking maps
	p.lastUsedMu.Lock()
	p.lastUsed = make(map[*Browser]time.Time)
	p.lastUsedMu.Unlock()
}

// cleanupIdleBrowsers runs periodically to close browsers that have been idle too long
func (p *BrowserPoolWithTimeout) cleanupIdleBrowsers() {
	defer p.wg.Done()

	ticker := time.NewTicker(p.idleTimeout / 2) // Check twice as often as timeout
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.performCleanup()
		case <-p.stopCleanup:
			return
		}
	}
}

// performCleanup checks for idle browsers and replaces them with fresh ones
func (p *BrowserPoolWithTimeout) performCleanup() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return
	}

	now := time.Now()
	browsersToReplace := []*Browser{}

	// Check browsers currently in the pool
	poolSize := len(p.Pool)
	for i := 0; i < poolSize; i++ {
		select {
		case browser := <-p.Pool:
			p.lastUsedMu.RLock()
			lastUsed, exists := p.lastUsed[browser]
			p.lastUsedMu.RUnlock()

			if exists && now.Sub(lastUsed) > p.idleTimeout {
				browsersToReplace = append(browsersToReplace, browser)
			} else {
				// Put the browser back if it's not idle
				select {
				case p.Pool <- browser:
				default:
					// Pool is full, close the browser
					browser.Close()
					p.lastUsedMu.Lock()
					delete(p.lastUsed, browser)
					p.lastUsedMu.Unlock()
				}
			}
		default:
			// No more browsers in pool
			break
		}
	}

	// Replace idle browsers with fresh ones
	for _, browser := range browsersToReplace {
		browser.Close()
		p.lastUsedMu.Lock()
		delete(p.lastUsed, browser)
		p.lastUsedMu.Unlock()

		// Create a new browser and add it to the pool
		newBrowser := NewBrowser()
		select {
		case p.Pool <- newBrowser:
			p.lastUsedMu.Lock()
			p.lastUsed[newBrowser] = now
			p.lastUsedMu.Unlock()
		default:
			// Pool is full, close the new browser
			newBrowser.Close()
		}
	}
}
