// Package main
// dstats : displays hardware information on conky
// complete package: main.go memory.go discs.go rubbish.go cpu.go utils.go + my.PrintCommaInt
package main

import (
	"log"
	"os/user"
)

// Width of main conky
const lineWidth = 40

var homeDir string

func main() {
	// find homeDir (for discs.go & rubbish.go)
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	homeDir = usr.HomeDir

	memoryInfo()
	discInfo()
	rubbishBins()
	temperatures()
}
