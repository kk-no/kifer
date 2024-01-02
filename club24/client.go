package club24

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/kk-no/kifer/internal/conv"
	"github.com/kk-no/kifer/internal/files"
	"github.com/kk-no/kifer/internal/input"
	"github.com/sclevine/agouti"
)

type Client struct {
	dir  string
	user string
	pass string

	// Used when temporary Sleep is required, such as page rendering.
	waitTime time.Duration
}

func New(dir, user, pass string, wait time.Duration) *Client {
	return &Client{
		dir:      dir,
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

	MinMoveCount int
	MinRate      int
}

func (c *Client) Login(page *agouti.Page) error {
	if err := c.access(page); err != nil {
		return fmt.Errorf("access failed: %w", err)
	}
	if err := c.login(page); err != nil {
		return fmt.Errorf("login failed: %w", err)
	}
	if !c.isLogin(page) {
		// NOTE: Login may not be successful with reCAPTCHA.
		//  Manual login required.
		fmt.Print("Please press the Enter key after manually logging in...")
		input.WaitEnter(30 * time.Second)
		time.Sleep(c.waitTime)
	}
	return nil
}

func (c *Client) Download(page *agouti.Page, config *DownloadConfig) error {
	if err := c.fillAndSearch(page, config); err != nil {
		return err
	}
	if err := c.download(page); err != nil {
		return err
	}
	if err := c.unzipFilter(config.MinMoveCount, config.MinRate); err != nil {
		return err
	}
	return nil
}

func (c *Client) access(page *agouti.Page) error {
	if err := page.Navigate(kifuPageURL); err != nil {
		return err
	}
	time.Sleep(c.waitTime)
	return nil
}

func (c *Client) login(page *agouti.Page) error {
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
	return nil
}

func (c *Client) isLogin(page *agouti.Page) bool {
	if _, err := page.FindByXPath(kifuUser1FormXPath).Active(); err != nil {
		return false
	}
	return true
}

func (c *Client) fillAndSearch(page *agouti.Page, config *DownloadConfig) error {
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
	return nil
}

func (c *Client) download(page *agouti.Page) error {
	if err := page.FindByXPath(kifuDownloadButtonXPath).Click(); err != nil {
		return err
	}
	time.Sleep(c.waitTime * 5)
	return nil
}

func FilterByMoveCountAndRate(minMove, minRate int) files.FilterFunc {
	return func(r io.Reader) bool {
		rate1, rate2, count := 0, 0, 0
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			count++
			switch count {
			case 5: // Extract `9999` from `先手：user1(9999)`
				rate1 = extractRating(scanner.Text())
			case 6: // Extract `9999` from `後手：user2(9999)`
				rate2 = extractRating(scanner.Text())
			}
		}
		if err := scanner.Err(); err != nil {
			return false
		}
		return minMove+fileMetadataCount < count && minRate < rate1 && minRate < rate2
	}
}

func (c *Client) unzipFilter(minMove, minRate int) error {
	return files.UnzipInDirectory(c.dir, FilterByMoveCountAndRate(minMove, minRate))
}

var ratingExtractFormat = regexp.MustCompile(`\(([0-9]+)\)`)

func extractRating(s string) int {
	match := ratingExtractFormat.FindStringSubmatch(s)
	if len(match) != 2 {
		return 0
	}
	return conv.Atoi(match[1])
}
