/*****************************************
* mpdshow                                *
*	dependencies: mpc mutagen-inspect    *
******************************************/
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	log "github.com/subchen/go-log"
	"github.com/subchen/go-log/formatters"
)

const musicFolder = "/tmp/mpd/"
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
			mpc("stop")          // quit & stop playing
			run("mpd", "--kill") // kill the server
			app.Stop()
		} else if event.Rune() == 'q' { // quit but still play
			app.Stop()
		} else if event.Rune() == 'r' { // redraw display
			dirty = true
		} else if event.Rune() == '[' { // previous track
			previousTrack()
		} else if event.Rune() == '<' { // backwards 10 minutes
			mpc("-q", "seek", "-600")
		} else if event.Rune() == ',' { // backwards 1 minute
			mpc("-q", "seek", "-60")
		} else if event.Key() == tcell.KeyLeft { // backwards 10 seconds
			mpc("-q", "seek", "-10")
		} else if event.Rune() == ' ' { // toggle pause
			mpc("-q", "toggle")
		} else if event.Key() == tcell.KeyRight { // forwards 10 seconds
			mpc("-q", "seek", "+10")
		} else if event.Rune() == '.' { // forwards 1 minute
			mpc("-q", "seek", "+60")
		} else if event.Rune() == '>' { // forwards 10 minutes
			mpc("-q", "seek", "+600")
		} else if event.Rune() == ']' { // next track
			nextTrack()
			// } else if event.Rune() == 's' { // toggle shuffle
			// 	mpc("-q", "shuffle")

			// change source
		} else if event.Key() == tcell.KeyBackspace2 { // play my mix
			setMusic("/share/music")
			mpc("-q", "--wait", "load", "my")
			mpc("-q", "--wait", "shuffle")
			mpc("-q", "play")
		} else if event.Rune() == '1' { // BBC World Service
			radioStation("http://bbcwssc.ic.llnwd.net/stream/bbcwssc_mp1_ws-eieuk")
		} else if event.Rune() == '2' { // Ancient FM
			radioStation("http://5.152.208.98:8058")
		} else if event.Rune() == '3' { // 181.Fm
			radioStation("http://relay.181.fm:8060")
		} else if event.Rune() == '4' { // Radionetz.de
			radioStation("http://0n-oldies.radionetz.de/0n-oldies.aac")
		}
		return event
	})

	// are there command line arguments?
	args := os.Args
	if len(args) > 1 {
		setMusic(args[1:]...)
		mpc("-q", "add", "/")
		mpc("-q", "play")
	}

	go updateAll()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func updateAll() {
	for {
		time.Sleep(refreshInterval)
		if err := mpc("-f", "%file%\n%title%"); err == nil {
			if currentFile.name != mpcArray[mpcFilename] && len(mpcArray[mpcFilename]) > 0 {
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
		statusBar.SetText("[yellow]" + str + progressBar(width-len(str)))
	} else {
		statusBar.SetText(pprintTime(currentFile.jetzt))
	}
}

func init() {
	log.Default.Level = log.WARN
	log.Default.Formatter = new(formatters.TextFormatter)

	dirty = true
	status, _ := run("pidof", "mpd")
	if status == nil { // not running: start it
		run("mpd")
	}
}

func nextTrack() {
	if cuesheet == nil {
		mpc("-q", "next")
	} else {
		nextTrack := cuesheet.getNextCue()
		if nextTrack != nil {
			mpc("-q", "seek", nextTrack.StarttimeInx)
		}
	}
}
func previousTrack() {
	if cuesheet == nil {
		mpc("-q", "prev")
	} else {
		previousTrack := cuesheet.getPreviousCue()
		mpc("-q", "seek", previousTrack.StarttimeInx)
	}
}

func radioStation(station string) {
	mpc("-q", "clear")
	mpc("-q", "add", station)
	mpc("-q", "--wait", "update")
	mpc("-q", "play")

}

func setMusic(paths ...string) {
	err := os.RemoveAll(musicFolder)
	if err != nil {
		log.Errorf("Cannot delete %s\n", musicFolder)
	}
	err = os.Mkdir(musicFolder, 0755)
	if err != nil {
		log.Errorf("Cannot create %s\n", musicFolder)
	}
	for _, path := range paths {
		abspath, err := filepath.Abs(path)
		if err != nil {
			log.Errorf("Error finding absolute path for %s\n", path)
		}
		err = os.Symlink(abspath, musicFolder+filepath.Base(path))
		if err != nil {
			log.Errorf("Cannot link %s to %s", abspath, musicFolder)
		}
	}
	mpc("-q", "clear")
	mpc("-q", "--wait", "update")
}
