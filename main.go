package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"time"

	"github.com/kk-no/kifer/club24"
	"github.com/kk-no/kifer/internal/dirs"
	"github.com/sclevine/agouti"
)

const (
	year  = 2023
	month = 12

	dateRange = 2
	minMove   = 50
	minRate   = 2700
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Print("failed to get current user")
		return
	}
	kifDir, err := dirs.CreateDirIfNotExist(fmt.Sprintf("%s/Downloads/kif", usr.HomeDir))
	if err != nil {
		log.Printf("failed to create and check directory: %s", err)
		return
	}
	driver := agouti.ChromeDriver(
		agouti.Browser("chrome"),
		agouti.ChromeOptions("prefs", map[string]interface{}{
			"download.default_directory": kifDir,
		}))
	if err := driver.Start(); err != nil {
		log.Println("failed to driver start", err)
		return
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		log.Printf("failed to open new page: %s", err)
		return
	}
	defer page.CloseWindow()

	cli := club24.New(kifDir, os.Getenv("CLUB_USER"), os.Getenv("CLUB_PASS"), 5*time.Second)
	if err := cli.Login(page); err != nil {
		log.Printf("failed to login: %s", err)
		return
	}

	start := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 0, dateRange)

	// Download in specified month.
	for {
		dlConf := &club24.DownloadConfig{
			User1:        "Hefeweizen",
			User2:        "",
			Start:        start,
			End:          end,
			MinMoveCount: minMove,
			MinRate:      minRate,
		}
		if err := cli.Download(page, dlConf); err != nil {
			log.Printf("failed to download from club 24: %s", err)
			break
		}
		start = end.AddDate(0, 0, 1)
		if start.Year() != year || start.Month() != month {
			break
		}
		end = start.AddDate(0, 0, dateRange)
		if end.Year() != year || end.Month() != month {
			// specify the end of the month.
			end = time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local).AddDate(0, 0, -1)
		}
	}
}
