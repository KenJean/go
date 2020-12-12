/*****************************************
* mocptracks                             *
*	dependencies: mocp mutagen-inspect   *
******************************************/
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	log "github.com/subchen/go-log"
	"github.com/subchen/go-log/formatters"
)

const refreshInterval = 1000 * time.Millisecond

var (
	app       *tview.Application
	textBox   *tview.TextView
	statusBar *tview.TextView
	dirty     bool
)

func main() {
	// set-up GUI
	app = tview.NewApplication()
	textBox = tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter).SetWordWrap(true)
	textBox.SetBackgroundColor(tcell.NewHexColor(0x000087))
	statusBar = tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	statusBar.SetBackgroundColor(tcell.NewHexColor(0x5f0000))
	flex := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(textBox, 0, 1, false).AddItem(statusBar, 1, 1, false)

	// Map keys
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// player control
		if event.Key() == tcell.KeyCtrlQ {
			mocp("-x") //kill mocp server
			app.Stop()
		} else if event.Rune() == 'q' { // quit without nuking server
			app.Stop()
		} else if event.Rune() == 'r' { // redraw display
			dirty = true
		} else if event.Rune() == '[' { // previous track
			previousTrack()
		} else if event.Rune() == '<' { // backwards 10 minutes
			mocp("-k", "-600")
		} else if event.Rune() == ',' { // backwards 1 minute
			mocp("-k", "-60")
		} else if event.Key() == tcell.KeyLeft { // backwards 10 seconds
			mocp("-k", "-10")
		} else if event.Rune() == ' ' { // toggle pause
			mocp("-G")
		} else if event.Key() == tcell.KeyRight { // forwards 10 seconds
			mocp("-k", "+10")
		} else if event.Rune() == '.' { // forwards 1 minute
			mocp("-k", "+60")
		} else if event.Rune() == '>' { // forwards 10 minutes
			mocp("-k", "+600")
		} else if event.Rune() == ']' { // next track
			nextTrack()
		} else if event.Key() == tcell.KeyUp { // louder
			mocp("-v", "+5")
		} else if event.Key() == tcell.KeyDown { // softer
			mocp("-v", "-7")
		} else if event.Rune() == 's' { // toggle shuffle
			mocp("-t", "shuffle")

			// change source
		} else if event.Rune() == 'p' { // my mix
			mocp("-p")
		} else if event.Rune() == '1' { // BBC
			mocp("-l", "http://bbcwssc.ic.llnwd.net/stream/bbcwssc_mp1_ws-eieuk")
		} else if event.Rune() == '2' { // BBC
			mocp("-l", "http://5.152.208.98:8058")
		} else if event.Rune() == '3' { // BBC
			mocp("-l", "http://relay.181.fm:8060")
		} else if event.Rune() == '4' { // BBC
			mocp("-l", "http://0n-oldies.radionetz.de/0n-oldies.aac")
		}
		return event
	})

	// are there command line arguments?
	args := os.Args
	if len(args) > 1 {
		// start server (disable shuffle)
		run("mocp", "-O", "Shuffle=no", "-S")
		time.Sleep(1000 * time.Millisecond)
		var parameters []string
		parameters = append(parameters, "-l")
		parameters = append(parameters, args[1:]...)
		run("mocp", parameters...)
	}

	go updateAll()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func updateAll() {
	for {
		time.Sleep(refreshInterval)
		if err := mocp("-i"); err == nil {
			if currentFile.name != mocpArray[mocpFilename] && len(mocpArray[mocpFilename]) > 0 {
				currentFile = newFileObject()
			}
			currentFile.update()
			app.QueueUpdateDraw(func() {
				tutti()
			})
		}
	}
}

// display volles Werk
// calls track.pprint() to print track info
func tutti() {
	if dirty {
		currentTrack.pprint()
	}
	// status bar: All time-related data
	if currentFile.dauer > 0 {
		str := ""
		if currentTrack.TrackNo > 0 {
			if currentTrack.TracksTotal > 0 {
				str = fmt.Sprintf("%2d ∈ %-2d ", currentTrack.TrackNo, currentTrack.TracksTotal)
			} else {
				str = fmt.Sprintf("   %2d   ", currentTrack.TrackNo)
			}
		} else {
			str = "        "
		}
		str += fmt.Sprintf(" %s ∑ %s  ", pprintTime(currentFile.jetzt), pprintTime(currentFile.dauer-currentFile.jetzt))
		_, _, width, _ := statusBar.GetInnerRect()
		statusBar.SetText("[yellow]" + str + progressBar(width-len(str), currentFile.jetzt*100/currentFile.dauer))
	} else {
		statusBar.SetText(pprintTime(currentFile.jetzt))
	}
}

func init() {
	log.Default.Level = log.WARN
	log.Default.Formatter = new(formatters.TextFormatter)

	dirty = true
}

func nextTrack() {
	if cuesheet == nil {
		mocp("-f")
	} else {
		nextTrack := cuesheet.getNextCue()
		if nextTrack != nil {
			mocp("-j", fmt.Sprintf("%ds", nextTrack.Starttime))
		}
	}
}
func previousTrack() {
	if cuesheet == nil {
		mocp("-r")
	} else {
		previousTrack := cuesheet.getPreviousCue()
		mocp("-j", fmt.Sprintf("%ds", previousTrack.Starttime))
	}
}
