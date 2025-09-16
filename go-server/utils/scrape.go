package utils

import (
	"context"

	"github.com/chromedp/chromedp"
)

type BrowserScraper struct {
	Ctx    context.Context
	cancel context.CancelFunc
}

func NewBrowserScraper() *BrowserScraper {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.Flag("user-data-dir", "/tmp/chromedp-profile"), // persistent profile dir
		chromedp.Flag("disable-site-isolation-trials", true),    // fewer processes (less isolation)
		chromedp.Flag("js-flags", "--max-old-space-size=96"),    // shrink V8 heaps (risky if pages are big)
		chromedp.Flag("disable-dev-shm-usage", true),            // avoid tiny /dev/shm on small VMs
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("window-size", "1280,800"),
		chromedp.Flag("mute-audio", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx)

	return &BrowserScraper{
		Ctx:    ctx,
		cancel: cancel,
	}
}

func (bs *BrowserScraper) Close() {
	bs.cancel()
}
