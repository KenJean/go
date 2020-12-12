package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	log "github.com/subchen/go-log"
)

type qObject struct {
	// sheet is 1-based
	tracks map[int]*trackObject
}

var (
	cuesheet *qObject
)

func newQObject() *qObject {
	qq := new(qObject)
	qq.tracks = make(map[int]*trackObject)
	return qq
}

func getExternalCuesheet(soundFilename string) (*qObject, error) {
	inx := strings.LastIndexByte(soundFilename, '.')
	cueFilename := soundFilename[:inx] + ".cue"
	if _, err := os.Stat(cueFilename); os.IsNotExist(err) {
		// file does not exist
		return nil, nil
	}
	file, err := os.Open(cueFilename)
	if err != nil {
		log.Warnln("cue file exists, but unable to be opened:", err)
		return nil, err
	}
	defer file.Close()

	var cueText []string
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Warnln("Error reading:, err")
				return nil, err
			}
		}
		cueText = append(cueText, strings.TrimSpace(line))
	}
	return makeCuesheet(cueText), nil
}

func makeCuesheet(cueText []string) *qObject {
	cueObject := newQObject()
	trackno := 0
	var cue *trackObject
	for _, txt := range cueText {
		// parse track tags in the form of: cue_track<#>_<tag> =value
		if strings.Contains(txt, "ue_track") {
			chop := strings.Split(txt, "_")
			data := strings.Split(chop[2], "=")
			trackNo, _ := strconv.Atoi(chop[1][5:])
			var track *trackObject
			if cueObject.tracks[trackNo] == nil {
				track = &trackObject{}
			} else {
				track = cueObject.tracks[trackNo]
			}
			tag := strings.ToUpper(data[0])
			switch tag {
			case "COMPOSER":
				track.Composer = data[1]
			case "ARTIST":
				track.Artist = data[1]
			case "BAND":
				track.Band = data[1]
			case "CONDUCTOR":
				track.Conductor = data[1]
			case "COMMENT":
				track.Comment = data[1]
			case "DATE":
				track.Date = data[1]
			}
			cueObject.tracks[trackNo] = track
			continue
		}
		// Skip over lines that do not apply to individual tracks (i.e. before trackno is set)
		if trackno == 0 && strings.HasPrefix(txt, "TRACK") == false && strings.Contains(txt, "AUDIO") == false {
			continue
		}

		// format of track information in cuesheet:
		// TRACK <#> AUDIO
		// TITLE "title"
		// PERFORMER "last,first, last,first …" or "last,first / last,first …"
		// INDEX 01 mm:ss:ff   ← this is always the last line
		if strings.HasPrefix(txt, "TRACK") {
			words := strings.Fields(txt)
			trackno, _ = strconv.Atoi(words[1])
			if cueObject.tracks[trackno] == nil {
				cue = &trackObject{}
			} else {
				cue = cueObject.tracks[trackno]
			}
		} else if strings.HasPrefix(txt, "PERFORMER") {
			words := strings.Split(txt, "\"")
			cue.Artist = words[1]
		} else if strings.HasPrefix(txt, "TITLE") {
			words := strings.Split(txt, "\"")
			cue.Title = words[1]
		} else if strings.HasPrefix(txt, "INDEX 01") { // this is the last line
			words := strings.Split(txt, "INDEX 01")
			times := strings.Split(words[1], ":")
			mmStr := strings.TrimSpace(times[0])
			ssStr := strings.TrimSpace(times[1])
			cue.StarttimeInx = mmStr + ":" + ssStr
			mm, _ := strconv.Atoi(mmStr)
			ss, _ := strconv.Atoi(ssStr)
			cue.Starttime = mm*60 + ss
			// as this is the last line, complete the cue and place in cuesheet
			cue.TrackNo = trackno
			cueObject.tracks[trackno] = cue
		}
	}
	return cueObject
}

func (q *qObject) getCueFromTime() *trackObject {
	var newCue *trackObject
	var trackNo int
	for trackNo = 0; trackNo < len(q.tracks); trackNo++ {
		cue := q.tracks[trackNo+1] // cuesheet is 1-based
		if cue.Starttime > currentFile.jetzt {
			break
		} else {
			newCue = cue
		}
	}
	if currentTrack.TrackNo != trackNo {
		dirty = true
	}
	if newCue != nil {
		track := new(trackObject)
		track.union(*newCue)
		return track
	}
	return currentTrack
}

func (q *qObject) getNextCue() *trackObject {
	nextQ := currentTrack.TrackNo + 1
	if nextQ > len(q.tracks) {
		return nil
	}
	return q.tracks[nextQ]
}

func (q *qObject) getPreviousCue() *trackObject {
	// Go to trackNo -1 *only if* within first 3 seconds of start
	// this to manually compesate for inaccuracies in decoding compressed audio
	if currentTrack.TrackNo > 1 && currentFile.jetzt-currentTrack.Starttime < 3 {
		previousQ := currentTrack.TrackNo - 1
		return q.tracks[previousQ]
	}
	return currentTrack
}
