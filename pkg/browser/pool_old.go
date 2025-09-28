package browser

func Create(size int) *BrowserPool {
	pool := make(chan *Browser, size)
	for i := range size{
		pool <- NewBrowser()
	}
}

func (browsers *BrowserPool) Accuire() *Browser{
	return <- browsers.Pool
}

func (browsers *BrowserPool) Release(browser *Browser) {
	browsers.Pool <- browser	
}