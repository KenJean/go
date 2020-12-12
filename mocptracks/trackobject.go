package main

import (
	"reflect"
	"strings"
)

type trackObject struct {
	// should be constructed with ints & strings only for union() to work
	// names must be "exportable" (i.e. begin with caps)
	Album       string
	TrackNo     int
	TracksTotal int
	Composer    string
	Title       string
	Artist      string
	Band        string
	Conductor   string
	Comment     string
	Date        string
	Label       string
	Starttime   int
}

var (
	currentTrack *trackObject
)

func (t trackObject) pprint() {
	var displayText []string
	if len(t.Album) > 0 {
		displayText = append(displayText, "[red]»[orange] "+t.Album+" [red]«")
	}

	if len(t.Composer) > 0 {
		names := parseArtists(t.Composer, "-")
		// Multiple composers will already be combined into the first slice
		displayText = append(displayText, "[fuchsia]"+names[0])
	}

	displayText = append(displayText, "[white]"+t.Title)

	if len(t.Artist) > 0 {
		artists := parseArtists(t.Artist, " • ")
		for _, line := range artists {
			displayText = append(displayText, "[lawngreen]"+line)
		}
	}

	if len(t.Band) > 0 {
		bands := parseBands(t.Band)
		for _, line := range bands {
			displayText = append(displayText, "[aqua]"+line)
		}
	}

	if len(t.Conductor) > 0 {
		displayText = append(displayText, "[crimson]——— "+firstNameFirst(t.Conductor)+" ———")
	}

	if len(t.Comment) > 0 {
		displayText = append(displayText, "[grey]"+t.Comment)
	}
	if len(t.Date) > 0 {
		displayText = append(displayText, "[darkcyan]"+t.Date)
	}
	if len(t.Label) > 0 {
		displayText = append(displayText, "[chocolate]"+t.Label)
	}

	// Put whole thing into the TextView
	textBox.SetText(strings.Join(displayText, "\n"))
	dirty = false
}

// union of two trackObjects
// only copy from t2 if field is not empty
// struct can only be constructed with ints & strings
func (t *trackObject) union(t2 trackObject) {
	// if len(t2.album) > 0 {
	// 	t.album = t2.album
	// } // etc.
	ss := reflect.ValueOf(t).Elem()
	ss2 := reflect.ValueOf(&t2).Elem()
	for ii := 0; ii < ss.NumField(); ii++ {
		ff := ss.Field(ii)
		ff2 := ss2.Field(ii)
		if ff != ff2 {
			vv2 := reflect.Value(ff2)
			// only replace ff with ff2 if ff2 is not empty
			if vv2.Kind() == reflect.Int && vv2.Int() > 0 {
				ff.Set(reflect.Value(ff2))
			} else if vv2.Kind() == reflect.String && vv2.Len() > 0 {
				ff.Set(reflect.Value(ff2))
			}
		}
	}
}

// parse list of artists in the form of "lastname,firstname/lastname,firstname …"
// does proper wrapping
// return a splice of names with first name first
func parseArtists(str, delimiter string) []string {
	var result, separated []string
	// first divide names into a splice
	if strings.Contains(str, ", ") {
		separated = strings.Split(str, ", ")
	} else {
		separated = strings.Split(str, "/")
	}
	// loop over splice and rearrange names to be first name first
	for _, lastFirst := range separated {
		firstLast := firstNameFirst(strings.TrimSpace(lastFirst))
		result = append(result, firstLast)
	}
	// manual wrapping
	_, _, width, _ := textBox.GetInnerRect()
	result = wrap(result, width, delimiter)
	return result
}

// parse list of bands delimited by "/"
// does proper wrapping
func parseBands(str string) []string {
	var result []string
	// first divide names into a splice
	separated := strings.Split(str, "/")
	// trim each entry
	for _, band := range separated {
		result = append(result, strings.TrimSpace(band))
	}
	// manual wrapping
	_, _, width, _ := textBox.GetInnerRect()
	result = wrap(result, width, " • ")
	return result
}

// Take Last-name,First-name and return First-name Last-name
func firstNameFirst(name string) string {
	result := name // in case already first name first
	firstLast := strings.Split(name, ",")
	if len(firstLast) > 1 {
		result = strings.TrimSpace(firstLast[1]) + " " + strings.TrimSpace(firstLast[0])
	}
	return result
}

// wrap the text using the slices in names as the atomic unit
// adds delimiter to separate the names but removes it from the end of a line
func wrap(names []string, boxWidth int, delimiter string) []string {
	var out []string
	var line string
	for inx, name := range names {
		_line := line + name + delimiter
		// use len() instead of utf8.RuneCountInString();
		// inx > 0: don't wrap if first item
		if len(_line) > boxWidth && inx > 0 {
			out = append(out, strings.TrimSpace(strings.TrimRight(line, delimiter)))
			line = name + delimiter
			continue
		}
		line += name + delimiter
	}
	out = append(out, strings.TrimRight(line, delimiter))
	return out
}

func init() {
	currentTrack = new(trackObject)
}
