package main

import (
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/subchen/go-log"
	"golang.org/x/text/encoding/charmap"
)

type fileObject struct {
	// fields to be retrieved from mocp -i
	name  string
	jetzt int
	dauer int
}

var (
	currentFile *fileObject
)

// make a new fileObject from the output of mocp -i (map of string)
func newFileObject() *fileObject {
	ptr := new(fileObject)
	dirty = true
	ptr.name = mocpArray[mocpFilename]
	ptr.jetzt, _ = strconv.Atoi(mocpArray[mocpCurrentsec])
	ptr.dauer, _ = strconv.Atoi(mocpArray[mocpTotalsec])

	var cuesheetRaw []string
	if strings.HasPrefix(ptr.name, "http:") {
		currentTrack = podcast(ptr.name)
		var err error
		currentTrack.Title, err = charmap.Windows1250.NewDecoder().String(mocpArray[mocpSongtitle])
		if err != nil {
			log.Warnln("unable to convert title to utf8")
		}
	} else {
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
	f.jetzt, _ = strconv.Atoi(mocpArray[mocpCurrentsec])
	if cuesheet != nil {
		track := cuesheet.getCueFromTime()
		currentTrack.union(*track)
	} else {
		currentTrack.Title = mocpArray[mocpSongtitle]
		dirty = true
	}
}

func init() {
	currentFile = new(fileObject)
}
