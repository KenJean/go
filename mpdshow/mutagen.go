package main

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// mutagen: parses filename to determine type of sound file and then triggers the proper subroutine
// returns: pointer to a trackObject & a cuesheet in raw text
func mutagen(filename string) (*trackObject, []string) {
	var track *trackObject
	mutagenOutput, cuesheetRaw := mutagenExec(filename)
	switch filepath.Ext(filename) {
	case ".ape":
		track = ape(mutagenOutput)
	case ".flac":
		track = flac(mutagenOutput)
	case ".m4a":
		track = m4a(mutagenOutput)
	case ".mp2":
		track = mp3(mutagenOutput)
	case ".mp3":
		track = mp3(mutagenOutput)
	case ".ogg":
		track = ogg(mutagenOutput)
	case ".opus":
		track = ogg(mutagenOutput)
	default:
		track = new(trackObject)
	}
	return track, cuesheetRaw
}

// mutagenExec: actual invocation of mutagen-inspect by the subroutines to filters out a cuesheet
// returns:	mutagenTxt = everything else
//			cuesheetTxt = plain text of what constitutes a proper cuesheet

func mutagenExec(filename string) ([]string, []string) {
	output, err := run("mutagen-inspect", filename)
	if err != nil {
		return nil, nil
	}
	var cuesheetTxt, mutagenTxt []string
	readingCuesheet := false
	for _, text := range output {
		if readingCuesheet {
			// If line contains "=" then it is no longer part of the cuesheet
			if strings.Contains(text, "=") {
				readingCuesheet = false
			} else {
				cuesheetTxt = append(cuesheetTxt, strings.TrimSpace(text))
				continue // NB. loop to next line
			}
		}
		if strings.Contains(text, "uesheet") {
			readingCuesheet = true
			cuesheetTxt = append(cuesheetTxt, text)
		} else if strings.Contains(text, "ue_") {
			cuesheetTxt = append(cuesheetTxt, text)
		} else {
			mutagenTxt = append(mutagenTxt, text)
		}
	}
	return mutagenTxt, cuesheetTxt
}

/* ***************************************************
   subroutines to process individual sound file types
******************************************************/

func ape(mutagen []string) *trackObject {
	track := new(trackObject)
	for _, line := range mutagen {
		if strings.Contains(line, "=") == false {
			continue
		}
		var fields []string
		fields = strings.SplitN(line, "=", 2) // Split only on first occurrence of '='
		key := strings.ToUpper(fields[0])
		value := strings.TrimSpace(fields[1])
		switch key {
		case "ALBUM":
			track.Album = value
		case "COMPOSER":
			track.Composer = value
		case "TITLE":
			track.Title = value
		case "ARTIST":
			track.Artist = value
		case "CONDUCTOR":
			track.Conductor = value
		case "BAND":
			track.Band = value
		case "COMMENT":
			track.Comment = value
		case "YEAR":
			track.Date = value
		case "DATE":
			track.Date = value
		case "LABEL":
			track.Label = value
		}
	}
	return track
}

func flac(mutagen []string) *trackObject {
	track := new(trackObject)
	var _composer, _artist, _band string
	for _, line := range mutagen {
		var fields []string
		if strings.Contains(line, "=") == false {
			continue
		}
		fields = strings.SplitN(line, "=", 2) // Split only on first occurrence of '='
		key := strings.ToUpper(fields[0])
		value := strings.TrimSpace(fields[1])
		switch key {
		case "ALBUM":
			track.Album = value
		case "COMPOSER":
			_composer += (value + " / ")
		case "TITLE":
			track.Title = value
		case "ARTIST":
			_artist += (value + " / ")
		case "CONDUCTOR":
			track.Conductor = value
		case "BAND":
			_band += (value + " / ")
		case "COMMENT":
			track.Comment = value
		case "YEAR":
			track.Date = value
		case "DATE":
			track.Date = value
		case "LABEL":
			track.Label = value
		}
	}
	track.Composer = strings.TrimSuffix(_composer, " / ")
	track.Artist = strings.TrimSuffix(_artist, " / ")
	track.Band = strings.TrimSuffix(_band, " / ")
	return track
}

func m4a(mutagen []string) *trackObject {
	track := new(trackObject)
	re := regexp.MustCompile("'(.*?)'")
	var _composer, _artist, _band string
	for _, line := range mutagen {
		var fields []string
		if strings.Contains(line, "=") == false {
			continue
		}
		fields = strings.SplitN(line, "=", 2) // Split only on first occurrence of '='
		key := strings.ToUpper(fields[0])
		value := strings.TrimSpace(fields[1])
		switch key {
		case "TRKN": // trkn=(track, totaltracks)
			text := value[1 : len(value)-1]
			fields := strings.Split(text, ", ")
			track.TrackNo, _ = strconv.Atoi(fields[0])
			if len(fields) == 2 {
				track.TracksTotal, _ = strconv.Atoi(fields[1])
			}
		case "ALBUM":
			track.Album = value
		case "©ALB":
			track.Album = value
		case "©WRT":
			_composer += (value + " / ")
		case "TITLE":
			track.Title = value
		case "©NAM":
			track.Title = value
		case "©ART":
			_artist += (value + " / ")
		case "©CON":
			track.Conductor = value
		case "----:COM.APPLE.ITUNES:CONDUCTOR":
			// ----:com.apple.iTunes:CONDUCTOR=MP4FreeForm('Nebolsin,Vasily', <AtomDataType.UTF8: 1>)
			conductorSlice := re.FindAllString(value, -1)
			track.Conductor = conductorSlice[0][1 : len(conductorSlice[0])-1]
		case "----:COM.APPLE.ITUNES:BAND":
			// ----:com.apple.iTunes:BAND=MP4FreeForm('Bolshoi Theatre', <AtomDataType.UTF8: 1>)
			bandSlice := re.FindAllString(value, -1)
			_band += (bandSlice[0][1:len(bandSlice[0])-1] + " / ")
		case "COMMENT":
			track.Comment = value
		case "©CMT":
			track.Comment = value
		case "©DAY":
			track.Date = value
		case "----:COM.APPLE.ITUNES:LABEL":
			// ----:com.apple.iTunes:LABEL=MP4FreeForm('YouTube', <AtomDataType.UTF8: 1>)
			labelSlice := re.FindAllString(value, -1)
			track.Label = labelSlice[0][1 : len(labelSlice[0])-1]
		}
	}
	track.Composer = strings.TrimSuffix(_composer, " / ")
	track.Artist = strings.TrimSuffix(_artist, " / ")
	track.Band = strings.TrimSuffix(_band, " / ")
	return track
}

func mp3(mutagen []string) *trackObject {
	track := new(trackObject)
	mutagenMap := make(map[string]string) // keys will be all upper-case
	for _, line := range mutagen {
		var fields []string
		if strings.Contains(line, "=") == false {
			continue
		}
		fields = strings.Split(line, "=") // Split only on first occurrence of '='
		if len(fields) == 2 {
			key := fields[0]
			mutagenMap[key] = strings.TrimSpace(fields[1])
		} else if len(fields) == 3 {
			if strings.HasPrefix(fields[0], "TXXX") { // BAND, DATE, LABEL
				key := strings.ToUpper(fields[1])
				mutagenMap[key] = strings.TrimSpace(fields[2])
			}
		} else if len(fields) == 4 { // COMMENT
			if fields[0] == "COMM" {
				mutagenMap[fields[0]] = strings.TrimSpace(fields[3])
			}
		}
	}
	for key, value := range mutagenMap { // keys are all upper-case
		switch key {
		case "TRCK": // either just a track number or in the form of "trackno/totaltracks"
			fields := strings.Split(value, "/")
			if len(fields) == 2 {
				track.TrackNo, _ = strconv.Atoi(fields[0])
				track.TracksTotal, _ = strconv.Atoi(fields[1])
			} else {
				track.TrackNo, _ = strconv.Atoi(value)
			}
		case "TALB":
			track.Album = value
		case "TCOM":
			track.Composer = value
		case "TIT2":
			track.Title = value
		case "TPE1":
			track.Artist = value
		case "TPE3":
			track.Conductor = value
		case "BAND":
			track.Band = value
		case "TPE2":
			track.Band = value
		case "COMM":
			if len(track.Comment) < len(value) {
				track.Comment = value
			}
		case "TDRC":
			if len(track.Date) < len(value) {
				track.Date = value
			}
		case "DATE":
			if len(track.Date) < len(value) {
				track.Date = value
			}
		case "LABEL":
			track.Label = value
		}
	}
	return track
}

func ogg(mutagen []string) *trackObject {
	track := new(trackObject)
	var _composer, _artist, _band string
	for _, line := range mutagen {
		var fields []string
		if strings.Contains(line, "=") == false {
			continue
		}
		fields = strings.SplitN(line, "=", 2) // Split only on first occurrence of '='
		key := strings.ToUpper(fields[0])
		value := strings.TrimSpace(fields[1])
		switch key {
		case "TRACKNUMBER":
			track.TrackNo, _ = strconv.Atoi(value)
		case "TRACKTOTAL":
			track.TracksTotal, _ = strconv.Atoi(value)
		case "ALBUM":
			track.Album = value
		case "COMPOSER":
			_composer += (value + " / ")
		case "TITLE":
			track.Title = value
		case "ARTIST":
			_artist += (value + " / ")
		case "CONDUCTOR":
			track.Conductor = value
		case "BAND":
			_band += (value + " / ")
		case "COMMENT":
			track.Comment = value
		case "DATE":
			track.Date = value
		case "LABEL":
			track.Label = value
		}
	}
	track.Composer = strings.TrimSuffix(_composer, " / ")
	track.Artist = strings.TrimSuffix(_artist, " / ")
	track.Band = strings.TrimSuffix(_band, " / ")
	return track
}
