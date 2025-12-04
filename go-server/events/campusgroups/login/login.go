// Package login manages autologin for campus groups,
// including headless browser actions and otp code retrieval from email.
// it responds with the minimum necessary response -- the auth cookie
package login

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

var SessionCookieName = "CG.SessionID"

type LoginParams struct {
	LoginEmail       string
	OtpEmail         string
	OtpEmailPassword string
}

func GetLoginCookie(loginURL string, p LoginParams) (string, error) {
	// create browser scraper
	bs := utils.NewBrowserScraper()
	defer bs.Close()

	ctx, cancel := chromedp.NewContext(bs.Ctx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// check if already logged in
	finalURL, err := navigateToLogin(ctx, loginURL)
	if err != nil {
		return "", fmt.Errorf("navigation failed: %w", err)
	}

	if finalURL != loginURL {
		// we got redirected, so we're already logged in.
		return getAuthCookie(ctx)
	}

	// create email client
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		return "", fmt.Errorf("failed to connect to email: %w", err)
	}
	defer func() {
		if err := c.Logout(); err != nil {
			// log, metric, etc. — don't return because we're in a defer
			log.Printf("imap logout failed: %v", err)
		}
	}()

	if err := c.Login(p.OtpEmail, p.OtpEmailPassword); err != nil {
		return "", fmt.Errorf("failed to login to email: %w", err)
	}

	// email snapshot + first stage of login flow
	baseline, err := triggerOtpAndSnapshot(c, ctx, p.LoginEmail, 10)
	if err != nil {
		return "", err
	}

	// grab the new otp
	otp, err := waitForNewOtp(c, baseline, 60*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to fetch OTP: %w", err)
	}

	// second stage
	if err := performLoginOtpStage(ctx, otp); err != nil {
		return "", fmt.Errorf("login otp stage failed: %w", err)
	}

	var postLoginURL string
	if err := chromedp.Run(ctx, chromedp.Location(&postLoginURL)); err != nil {
		return "", fmt.Errorf("failed to get post-login URL: %w", err)
	}

	return getAuthCookie(ctx)
}

func navigateToLogin(ctx context.Context, loginURL string) (string, error) {
	var finalURL string
	err := chromedp.Run(ctx,
		chromedp.Navigate(loginURL),

		chromedp.WaitReady("body", chromedp.ByQuery),

		chromedp.Location(&finalURL),
	)
	return finalURL, err
}

func performLoginEmailStage(ctx context.Context, email string) error {
	return chromedp.Run(ctx,
		// expand the hidden area
		chromedp.Click("#a-all-others-sign-in-below", chromedp.ByID),

		// wait for expansion
		chromedp.WaitVisible("#cornell_login", chromedp.ByID),

		chromedp.SendKeys("#login_email", email, chromedp.ByID),

		// ensure remember-me is ON
		chromedp.ActionFunc(func(ctx context.Context) error {
			var checked bool
			if err := chromedp.Evaluate(`document.querySelector("#remember_me").checked`, &checked).Do(ctx); err != nil {
				return err
			}
			if !checked {
				return chromedp.Click("#remember_me", chromedp.ByID).Do(ctx)
			}
			return nil
		}),

		chromedp.Click("#loginButton", chromedp.ByID),
	)
}

func triggerOtpAndSnapshot(c *client.Client, ctx context.Context, email string, lastN uint32) (map[uint32]*imap.Envelope, error) {
	baseline, err := fetchRecentEnvelopes(c, lastN)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch baseline emails: %w", err)
	}

	if err := performLoginEmailStage(ctx, email); err != nil {
		return nil, fmt.Errorf("login email stage failed: %w", err)
	}

	return baseline, nil
}

func performLoginOtpStage(ctx context.Context, otp string) error {
	return chromedp.Run(ctx,
		chromedp.WaitVisible("#otp_form", chromedp.ByID),

		chromedp.SendKeys("#otp", otp, chromedp.ByID),

		chromedp.Click("#otb_button", chromedp.ByID),

		chromedp.WaitNotPresent("#otp_form", chromedp.ByID),
	)
}

func getAuthCookie(ctx context.Context) (string, error) {
	var cookies []*network.Cookie
	if err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		cookies, err = network.GetCookies().Do(ctx)
		return err
	})); err != nil {
		return "", fmt.Errorf("failed to get cookies: %w", err)
	}

	for _, c := range cookies {
		if c.Name == SessionCookieName {
			return c.Value, nil
		}
	}

	return "", fmt.Errorf("cookie %s not found", SessionCookieName)
}
