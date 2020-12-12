package main

import (
	"path/filepath"
	"strings"

	log "github.com/subchen/go-log"
	"golang.org/x/text/encoding/charmap"
)

type fileObject struct {
	// fields to be retrieved from mcp
	name    string
	jetzt   int
	dauer   int
	percent int
}

var (
	currentFile *fileObject
)

// make a new fileObject from the output of mcp (string array)
func newFileObject() *fileObject {
	ptr := new(fileObject)
	dirty = true
	ptr.name = mpcArray[mpcFilename]
	ptr.jetzt = mpcTimings[mpcCurrentsec]
	ptr.dauer = mpcTimings[mpcTotalsec]
	ptr.percent = mpcTimings[mpcPercent]

	var cuesheetRaw []string
	if strings.Contains(ptr.name, "http:") {
		currentTrack = podcast(ptr.name)
		var err error
		currentTrack.Title, err = charmap.Windows1250.NewDecoder().String(mpcArray[mpcSongtitle])
		if err != nil {
			log.Warnln("unable to convert title to utf8")
		}
	} else {
		ptr.name = musicFolder + ptr.name
		currentTrack, cuesheetRaw = mutagen(ptr.name)
	}
	// Is there an embedded cuesheet?
	if len(cuesheetRaw) > 0 {
		cuesheet = makeCuesheet(cuesheetRaw)
	} else {
		var err error
		cuesheet, err = getExternalCuesheet(ptr.name)
		if err != nil {
			log.Warnln("getExternalCuesheet() error")
		}
	}
	if cuesheet != nil {
		currentTrack.TracksTotal = len(cuesheet.tracks)
	}
	// if all else fails (& not streaming), use file name as title
	if len(currentTrack.Title) == 0 && strings.HasPrefix(ptr.name, "http:") == false {
		currentTrack.Title = strings.TrimRight(filepath.Base(ptr.name), filepath.Ext(ptr.name))
	}
	return ptr
}

func (f *fileObject) update() {
	f.jetzt = mpcTimings[mpcCurrentsec]
	f.percent = mpcTimings[mpcPercent]
	if cuesheet != nil {
		track := cuesheet.getCueFromTime()
		currentTrack.union(*track)
	} else {
		currentTrack.Title = mpcArray[mpcSongtitle]
		dirty = true
	}
}

func init() {
	currentFile = new(fileObject)
}
