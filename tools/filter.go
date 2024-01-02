package main

import (
	"fmt"
	"log"
	"os/user"

	"github.com/kk-no/kifer/club24"
	"github.com/kk-no/kifer/internal/dirs"
	"github.com/kk-no/kifer/internal/files"
)

const (
	minMove = 50
	minRate = 2700
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
	if err := files.UnzipInDirectory(kifDir, club24.FilterByMoveCountAndRate(minMove, minRate)); err != nil {
		log.Printf("failed to unzip and filter in directory: %s", err)
		return
	}
}
