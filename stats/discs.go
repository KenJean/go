package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"my"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// partition : struct for disc partition info
type partition struct {
	name string
	path string
}

// findPartitions : Loops through /etc/mtab and filters out extraneous entries
//	returns a string slice of partition structs
func findPartitions() []partition {
	lines := getFromFile("/etc/mtab")
	tutti := []partition{}
	for _, line := range lines {
		strArray := strings.Fields(line)
		if len(strArray) > 0 {
			if strings.Contains(strArray[0], "/dev/sd") || (strArray[0] == "tmpfs" && strArray[1] == "/tmp") {
				if strArray[1] != "/netdrive" {
					part := partition{strArray[0], strArray[1]}
					tutti = append(tutti, part)
				}
			}
		}
	}
	return tutti
}

// getBaseline
//		label = name of the partition
//		size = amount of current free space
// returns:
//		size (unaltered) if this partition is new
//		baseline as read from persistence storage
func getBaseline(label string, size int) int {
	// set full path name for saving "size"
	iniDir := homeDir + "/.config/dstats/"
	var fullpath string
	if label == "/" {
		fullpath = iniDir + "root"
	} else {
		fullpath = iniDir + label
	}

	// read/write to disc
	content, err := ioutil.ReadFile(fullpath)
	if err != nil { // file does not (yet) exist
		err := os.MkdirAll(iniDir, 0777)
		if err != nil {
			log.Print(err)
		}
		err = ioutil.WriteFile(fullpath, []byte(strconv.Itoa(size)), 0644)
		if err != nil {
			log.Printf("Cannot open %v: %v", fullpath, err)
		}
	} else { // replace baseline with what is read from the disc
		size, _ = strconv.Atoi(string(content))
	}
	return size
}

// printPartitionInfo : path (to partition)
// prints out info for a single partition
func printPartitionInfo(path string) {
	var fiData syscall.Statfs_t
	syscall.Statfs(path, &fiData)
	// struct statvfs fiData;
	// 	unsigned long f_bsize    File name block size.
	// 	unsigned long f_frsize   Fundamental file name block size.
	// 	fsblkcnt_t    f_blocks   Total number of blocks on file name in units of f_frsize.
	// 	fsblkcnt_t    f_bfree    Total number of free blocks.
	// 	fsblkcnt_t    f_bavail   Number of free blocks available to non-privileged process.
	// 	fsfilcnt_t    f_files    Total number of file serial numbers.
	// 	fsfilcnt_t    f_ffree    Total number of free file serial numbers.
	// 	fsfilcnt_t    f_favail   Number of file serial numbers available to non-privileged process.
	// 	unsigned long f_fsid     File name ID.
	// 	unsigned long f_flag     Bit mask of f_flag values.
	// 	unsigned long f_namemax  Maximum filename length.

	if len(path) > 1 { // if not just '/'
		parsed := strings.Split(path, "/")
		path = parsed[len(parsed)-1]
	}

	free := int(fiData.Bsize) * int(fiData.Bavail)
	percentUsed := 100.0 * float64(fiData.Blocks-fiData.Bfree) / float64(fiData.Blocks-fiData.Bfree+fiData.Bavail)

	// Make persistence Object
	change := getBaseline(path, free) - free
	fmt.Printf("${color yellow}%-7.7s %15v ${color cyan}%16s\n${color 33FFFF}%3d%% ", path, my.PrintCommaInt(change), my.PrintCommaInt(free), int(percentUsed+0.5))
	printgraph(lineWidth-5, percentUsed)
}

// discInfo : entry point for this package
func discInfo() {
	partitions := findPartitions()
	for _, part := range partitions {
		printPartitionInfo(part.path)
	}
}
