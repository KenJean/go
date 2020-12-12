package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"my"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// exists : check if folder exists
//		should not be used for files
func exists(path string) bool {
	stat, err := os.Stat(path)
	if err == nil && stat.IsDir() {
		return true
	}
	return false
}

// folderSize : uses filepath.walk to loop through subfolders
//	NB. Does NOT sum folders
func folderSize(path string) int {
	var dirSize int
	readSize := func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			dirSize += int(file.Size())
		}
		return nil
	}
	filepath.Walk(path, readSize)
	return dirSize
}

// rubbishBins : sum the sizes of all that is currently in Trash
// 	Create a directory slice (rubbish) of all possible parent directories
func rubbishBins() {
	var rubbish []string // Rubbish will contain parent directories of Trash folders
	//    Each must end with "/"
	rubbish = append(rubbish, homeDir+"/.local/share/")
	rubbish = append(rubbish, "/tmp/")
	// Only if /share exists
	if exists("/share/") {
		rubbish = append(rubbish, "/share/")
	}
	// Add every device attached to /media/$USER
	usr, _ := user.Current()
	mountPoint := ""
	if exists("/media/") {
		mountPoint = "/media/" + usr.Username + "/"
	} else {
		mountPoint = "/run/media/" + usr.Username + "/"
	}
	if exists(mountPoint) {
		items, err := ioutil.ReadDir(mountPoint)
		if err != nil {
			log.Print(err)
		}
		for _, item := range items {
			if item.IsDir() {
				rubbish = append(rubbish, mountPoint+item.Name()+"/")
			}
		}
	}
	var total int
	// Loop through rubbish
	for _, folder := range rubbish {
		// Loop through each child directory, processing only those that contain the string "Trash"
		subfolders, err := ioutil.ReadDir(folder)
		if err != nil {
			log.Print(err)
		}
		for _, item := range subfolders {
			if item.IsDir() && strings.Contains(item.Name(), "Trash") {
				fullpath := folder + item.Name()
				total += folderSize(fullpath)
			}
		}
	}
	fmt.Printf("${color yellow}${alignc}Rubbish: %s\n", my.PrintCommaInt(total))
}
