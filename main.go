package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/kk-no/kifer/club24"
	"github.com/sclevine/agouti"
)

var now = time.Now()

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Print("failed to get current user")
		return
	}
	kifDir, err := createKifDirIfNotExist(fmt.Sprintf("%s/Downloads", usr.HomeDir))
	if err != nil {
		log.Println("failed to create and check directory")
		return
	}
	driver := agouti.ChromeDriver(agouti.Browser("chrome"), agouti.ChromeOptions("prefs", map[string]interface{}{
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

	cli := club24.New(os.Getenv("CLUB_USER"), os.Getenv("CLUB_PASS"), 5*time.Second)
	dlConf := &club24.DownloadConfig{
		User1: "Hefeweizen",
		User2: "",
		Start: time.Date(now.Year(), now.Month(), now.Day()-10, 0, 0, 0, 0, time.Local),
		End:   time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local),
	}
	if err := cli.Download(page, dlConf); err != nil {
		log.Printf("failed to download from club 24: %s\n", err)
		return
	}
	if err := unzipInDirectory(kifDir); err != nil {
		log.Println("failed to unzip and remove zip file")
		return
	}
}

func createKifDirIfNotExist(path string) (string, error) {
	kifDir := path + "/kif"
	if _, err := os.Stat(kifDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(kifDir, os.ModePerm); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return kifDir, nil
}

func unzipInDirectory(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".zip") {
			continue
		}
		if err := unzip(dir, file.Name()); err != nil {
			return err
		}
	}
	return nil
}

func unzip(dir, zipFile string) error {
	zipFilePath := filepath.Join(dir, zipFile)
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if err := createFile(dir, f); err != nil {
			return err
		}
	}
	if err := os.Remove(zipFilePath); err != nil {
		return err
	}
	return nil
}

func createFile(dir string, f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	out, err := os.OpenFile(filepath.Join(dir, f.Name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, rc); err != nil {
		return err
	}
	return nil
}
