package club24

import (
	"bufio"
	"fmt"
	"os"
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
	time.Sleep(c.waitTime)

	if err := page.FindByXPath(loginUserFormXPath).Fill(c.user); err != nil {
		return err
	}
	if err := page.FindByXPath(loginPassFormXPath).Fill(c.pass); err != nil {
		return err
	}
	if err := page.FindByXPath(loginSubmitButtonXPath).Click(); err != nil {
		return err
	}
	time.Sleep(c.waitTime)

	if _, err := page.FindByXPath(kifuUser1FormXPath).Active(); err != nil {
		// NOTE: Login may not be successful with reCAPTCHA.
		//  Manual login required.
		fmt.Print("Please press the Enter key after manually logging in...")
		for {
			if _, err := bufio.NewReader(os.Stdin).ReadString('\n'); err != nil {
				continue
			}
			break
		}
		fmt.Println("Resuming download")
		time.Sleep(c.waitTime)
	}

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
		if err := page.FindByName(kifuEndDateName).Fill(config.End.Format(dateOnly)); err != nil {
			return err
		}
	}
	if err := page.FindByXPath(kifuSearchButtonXPath).Click(); err != nil {
		return err
	}
	time.Sleep(c.waitTime)

	if err := page.FindByXPath(kifuDownloadButtonXPath).Click(); err != nil {
		return err
	}
	time.Sleep(c.waitTime * 5)
	return nil
}
