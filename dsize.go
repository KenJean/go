package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

var re = regexp.MustCompile("(\\d+)(\\d{3})")

func printCommaInt64(num int64) string {
	str := strconv.FormatInt(num, 10)
	for i := 0; i < (len(str)-1)/3; i++ {
		str = re.ReplaceAllString(str, "$1,$2")
	}
	return str
}

// folderSize : uses filepath.walk to loop through subfolders
//	NB. also sums folders
func folderSize(path string) int64 {
	var dirSize int64 = 0
	readSize := func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			dirSize += file.Size()
		}
		return nil
	}
	filepath.Walk(path, readSize)
	return dirSize
}

func main() {
	fmt.Println(printCommaInt64(folderSize(".")))
}
