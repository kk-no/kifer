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

var now = time.Now()

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Print("failed to get current user")
		return
	}
	kifDir, err := dirs.CreateDirIfNotExist(fmt.Sprintf("%s/Downloads/kif", usr.HomeDir))
	if err != nil {
		log.Println("failed to create and check directory")
		return
	}
	driver := agouti.ChromeDriver(
		agouti.Browser("chrome"),
		agouti.ChromeOptions("prefs", map[string]interface{}{
			"download.default_directory": kifDir,
		}))
	if err := driver.Start(); err != nil {
		log.Println("failed to driver start")
		return
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		log.Println("failed to open new page")
		return
	}
	defer page.CloseWindow()

	cli := club24.New(kifDir, os.Getenv("CLUB_USER"), os.Getenv("CLUB_PASS"), 5*time.Second)
	dlConf := &club24.DownloadConfig{
		User1:        "Hefeweizen",
		User2:        "",
		Start:        time.Date(now.Year(), now.Month(), now.Day()-10, 0, 0, 0, 0, time.Local),
		End:          time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local),
		MinMoveCount: 50,
		MinRate:      2700,
	}
	if err := cli.Download(page, dlConf); err != nil {
		log.Printf("failed to download from club 24: %s\n", err)
		return
	}
}
