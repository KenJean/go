package main

func podcast(url string) *trackObject {
	track := new(trackObject)
	switch url {
	case "http://bbcwssc.ic.llnwd.net/stream/bbcwssc_mp1_ws-eieuk":
		track.Album = "BBC World Service"
	case "http://5.152.208.98:8058":
		track.Album = "Ancient FM"
	case "http://relay.181.fm:8060":
		track.Album = "181.FM"
	case "http://0n-oldies.radionetz.de/0n-oldies.aac":
		track.Album = "Radionetz.de"
	}
	return track
}
