package main

import (
	"log"
	"os"
	"time"

	"github.com/kk-no/kifer/club24"
	"github.com/sclevine/agouti"
)

var now = time.Now()

func main() {
	driver := agouti.ChromeDriver(agouti.Browser("chrome"))
	if err := driver.Start(); err != nil {
		log.Panic("failed to driver start")
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		log.Panic("failed to open new page")
	}
	defer page.CloseWindow()

	cli := club24.New(os.Getenv("CLUB_USER"), os.Getenv("CLUB_PASS"), 3*time.Second)
	if err := cli.Download(page, &club24.DownloadConfig{
		User1: "Hefeweizen",
		User2: "",
		Start: time.Date(now.Year(), now.Month(), now.Day()-10, 0, 0, 0, 0, time.Local),
		End:   time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local),
	}); err != nil {
		log.Panicf("failed to download from club 24: %s\n", err)
	}
}
