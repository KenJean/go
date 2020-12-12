package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	log "github.com/subchen/go-log"
)

const (
	mpcFilename = iota
	mpcSongtitle
)
const (
	mpcTotalsec = iota
	mpcCurrentsec
	mpcPercent
)

var (
	mpcArray   [2]string
	mpcTimings [3]int
	numRe      *regexp.Regexp
)

// mpc(): runs mpc with parameters
// "mpc "-f %file%\n%title%" returns:
//		0: File name relative to musicFolder
// 		1: Title
//		2: [playing] #1/2 146:50/271:48 (54%)
//		3: volume: 50%   repeat: off   random: off   single: off   consume: off
//
// run("deadbeef", "--nowplaying-tf", "%path%¯%tracknumber%¯%playback_time%¯%playback_time_remaining%")
func mpc(parameters ...string) error {
	// if --quiet is set, don't parse (i.e. return immediately)
	if parameters[0] == "-q" {
		run("mpc", parameters...)
		return nil
	}
	output, err := run("mpc", parameters...)
	if err != nil {
		return err
	}
	if len(output) <= 3 { // currently not playing
		return err
	}
	for ii := range mpcArray {
		mpcArray[ii] = ""
	}
	for ii := range mpcTimings {
		mpcTimings[ii] = 0
	}
	mpcArray[mpcFilename] = output[0]  // filename
	mpcArray[mpcSongtitle] = output[1] // title
	// parse "[playing] #2/2 12:34/34:56 (12%)"
	fields := strings.Fields(output[2])
	timings := strings.SplitN(fields[2], "/", -1)
	mpcTimings[mpcTotalsec] = time2secs(timings[1])
	mpcTimings[mpcCurrentsec] = time2secs(timings[0])
	// parse "(12%)"
	result := numRe.FindAllString(fields[3], -1)
	mpcTimings[mpcPercent], _ = strconv.Atoi(result[0])

	return nil
}

// centre text within widget
// func centred(text string, width int) string {
// 	// find & nuke all '${color xxx}' substrings
// 	// re := regexp.MustCompile("[$][{].*?[}]")
// 	// plainText := re.ReplaceAllString(colourText, "")
// 	// plainTextLength := len(plainText)
// 	textLength := utf8.RuneCountInString(text)
// 	output := ""
// 	if textLength < width {
// 		pad := (((width*10 - textLength*10) / 2) + 5) / 10
// 		output = fmt.Sprint(strings.Repeat(" ", pad))
// 	}
// 	return output + text
// 	// fmt.Printf("%[1]*s", pad, colourText)
// }

// progressBar: prints a progress bar using the unicode █ & ░
func progressBar(width int) string {
	width -= 5 // length of " %3d%%"
	used := width * currentFile.percent / 100
	free := width - used
	return fmt.Sprintf("%s%s %3d%%", strings.Repeat("█", used), strings.Repeat("░", free), currentFile.percent)
}

// convert seconds to 00:00:00 and print out
func pprintTime(timeSecs int) string {
	ss := timeSecs % 60
	mm := timeSecs / 60
	hh := mm / 60
	mm %= 60
	return fmt.Sprintf("%d:%02d:%02d", hh, mm, ss)
}

// run an external program and capture its output
// returns lines in a slice
func run(programme string, parameters ...string) ([]string, error) {
	cmd := exec.Command(programme, parameters...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		log.Debugf("Error running %s\n", programme)
		return nil, err
	}
	return strings.Split(string(stdout.Bytes()), "\n"), nil
}

func time2secs(str string) int {
	fields := strings.SplitN(str, ":", -1)
	mins, _ := strconv.Atoi(fields[0])
	secs, _ := strconv.Atoi(fields[1])
	return mins*60 + secs
}

func init() {
	numRe = regexp.MustCompile("[0-9]+")
}
