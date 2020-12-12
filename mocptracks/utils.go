package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	log "github.com/subchen/go-log"
)

const (
	mocpFilename = iota
	mocpSongtitle
	mocpTotalsec
	mocpCurrentsec
)

var (
	mocpArray [4]string
)

// mocp(): runs mocp with parameters
// "mocp "-i" returns:
//       0: State: PLAY
//       1: File: (full path)
//       2: Title:
//       3: Artist:
//       4: SongTitle:
//       5: Album:
//       6: TotalTime:  226m
//       7: TimeLeft:  226m
//       8: TotalSec: 13576
//       9: CurrentTime: 00:14
//      10: CurrentSec: 14
//      11: Bitrate: 389kbps
//      12: AvgBitrate: 414kbps
//      13: Rate: 44kHz
// will make sure server is running if issuing a 'play' command
// run("deadbeef", "--nowplaying-tf", "%path%¯%tracknumber%¯%playback_time%¯%playback_time_remaining%")
func mocp(parameters ...string) error {
	// Check to see whether server is running (but only if invoking a play function)
	if parameters[0] == "-p" || parameters[0] == "-l" {
		status, _ := run("pidof", "mocp")
		if status == nil { // not running: start it
			run("mocp", "-S")
			time.Sleep(1000 * time.Millisecond)
		} else if parameters[0] == "-p" { // already running, so nuke it & restart it to re-enable default settings
			run("mocp", "-x")
			time.Sleep(1000 * time.Millisecond)
			run("mocp", "-S")
			time.Sleep(1000 * time.Millisecond)
		}
	}
	output, err := run("mocp", parameters...)
	if err != nil {
		return err
	}
	for ii := range mocpArray {
		mocpArray[ii] = ""
	}
	for _, line := range output {
		var fields []string
		fields = strings.SplitN(line, ":", 2) // Split only on first occurrence of ':'
		if len(fields[0]) > 0 {
			switch fields[0] {
			case "File":
				mocpArray[mocpFilename] = strings.TrimSpace(fields[1])
			case "SongTitle":
				mocpArray[mocpSongtitle] = strings.TrimSpace(fields[1])
			case "TotalSec":
				mocpArray[mocpTotalsec] = strings.TrimSpace(fields[1])
			case "CurrentSec":
				mocpArray[mocpCurrentsec] = strings.TrimSpace(fields[1])
			}
		}
	}
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
func progressBar(width, percent int) string {
	width -= 5 // length of " %3d%%"
	used := width * percent / 100
	free := width - used
	return fmt.Sprintf("%s%s %3d%%", strings.Repeat("█", used), strings.Repeat("░", free), percent)
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
