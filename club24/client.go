package club24

import (
	"time"

	"github.com/sclevine/agouti"
)

type Client struct {
	user string
	pass string

	// Used when temporary Sleep is required, such as page rendering.
	waitTime time.Duration
}

func New(user, pass string, wait time.Duration) *Client {
	return &Client{
		user:     user,
		pass:     pass,
		waitTime: wait,
	}
}

type DownloadConfig struct {
	User1 string
	User2 string
	Start time.Time
	End   time.Time
}

func (c *Client) Download(page *agouti.Page, config *DownloadConfig) error {
	if err := page.Navigate(kifuPageURL); err != nil {
		return err
	}
	c.sleep()

	if err := page.FindByXPath(loginUserFormXPath).Fill(c.user); err != nil {
		return err
	}
	if err := page.FindByXPath(loginPassFormXPath).Fill(c.pass); err != nil {
		return err
	}
	if err := page.FindByXPath(loginSubmitButtonXPath).Click(); err != nil {
		return err
	}
	c.sleep()

	if config.User1 != "" {
		if err := page.FindByXPath(kifuUser1FormXPath).Fill(config.User1); err != nil {
			return err
		}
	}
	if config.User2 != "" {
		if err := page.FindByXPath(kifuUser2FormXPath).Fill(config.User2); err != nil {
			return err
		}
	}
	if !config.Start.IsZero() {
		if err := page.FindByName(kifuStartDateName).Fill(config.Start.Format(dateOnly)); err != nil {
			return err
		}
	}
	if !config.End.IsZero() {
		if err := page.FindByXPath(kifuEndDateName).Fill(config.End.Format(dateOnly)); err != nil {
			return err
		}
	}
	if err := page.FindByXPath(kifuSearchButtonXPath).Click(); err != nil {
		return err
	}
	c.sleep()

	if err := page.FindByXPath(kifuDownloadButtonXPath).Click(); err != nil {
		return err
	}
	c.sleep()
	return nil
}

func (c *Client) sleep() {
	time.Sleep(c.waitTime)
}
