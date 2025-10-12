package browser

import (
	logger "consumer/pkg/log"
	"errors"
	"sync/atomic"
	"time"
)

func CreatePool(size int) (*BrowserPool, error) {
	pool := make(chan *Browser, size)
	dead := make(chan *Browser, size)

	for range size {
		browser := NewBrowser()
		browser.WarmUp()
		pool <- browser
	}

	active := atomic.Bool{}
	active.Store(true)
	
	browsers := &BrowserPool{
		size,
		pool,
		dead,
		&active,
	}
	
	if len(browsers.Pool) != browsers.Size {
		return nil, errors.New("browsers crashed") 
	}

	go browsers.Monitor()
	return browsers, nil
}

func (browsers *BrowserPool) Acquire() *Browser {
	for{
		if !browsers.Active.Load() {
			return nil
		}

		select {
			case browser := <-browsers.Pool:
				if browser == nil {
					continue
				}

				if !browser.IsAlive() {
					select {
						case browsers.Dead <- browser:
							logger.Warning("Browser found dead on acquire replacing... 💀")

						default:
							browser.Close()
					}
					continue
				}

				return browser

			case <-time.After(10 * time.Second):
				return nil
		}
	}
}

func (browsers *BrowserPool) Release(browser *Browser) {
	if browser == nil {
		return
	}

	if !browser.IsAlive() {
		select {
			case browsers.Dead <- browser:
				logger.Warning("Browser dead on release replacing... 💀")

			default:
				browser.Close()
		}

		return
	}

	select {
		case browsers.Pool <- browser:

		default:
			browser.Close()
	}
}

func (browsers *BrowserPool) Close() {
	browsers.Active.Swap(false)
	time.Sleep(100 * time.Millisecond)
	
	for {
		select {
			case browser := <-browsers.Pool:
				if browser != nil {
					browser.Close()
				}
			default:
				close(browsers.Pool)
				return
		}
	}
}

func (browsers *BrowserPool) Monitor(){
	for browsers.Active.Load() {
		select {
			case dead := <- browsers.Dead:
				if(dead != nil){
					dead.Close()
				}

				browser := NewBrowser()
				browser.WarmUp()

				select {
					case browsers.Pool <- browser:
						logger.Info("Replaced dead browser successfully ♻️")

					default:
						browser.Close()
				}
			
			case <-time.After(2 * time.Second):
		}
	}
}