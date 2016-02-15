package server

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	STARTED = iota
	COMPLETED
	STOPPED
)

type TorrentEvent struct {
	started   int
	completed int
	stopped   int
}

type TorrentRequestData struct {
	info_hash  string //20 byte sha1 hash
	peer_id    string //max len 20
	ip         string //optional
	port       int    // port number the peer is listening on
	uploaded   int    // base10 ascii amount uploaded so far
	downloaded int    // base10 ascii amount downloaded so far
	left       int    // # of bytes left to download (base 10 ascii)
	event      int
}

var ANNOUNCE_URL = "/announce"

func parseTorrentGetRequestURI(s string) map[string]interface{} {
	tmp := strings.Split(s, "?")
	tmp = strings.Split(tmp[1], "%26")
	result := make(map[string]interface{})
	for i := range tmp {
		if tmp[i] != ANNOUNCE_URL {
			data := strings.Split(tmp[i], "=")
			result[data[0]] = data[1]
		}
	}
	return result
}
func fillEmptyMapValues(torrentMap map[string]interface{}) *TorrentRequestData {
	_, ok := torrentMap["port"]
	if !ok {
		torrentMap["port"] = 0
	}
	_, ok = torrentMap["uploaded"]
	if !ok {
		torrentMap["uploaded"] = 0
	}
	_, ok = torrentMap["downloaded"]
	if !ok {
		torrentMap["downloaded"] = 0
	}
	_, ok = torrentMap["left"]
	if !ok {
		torrentMap["left"] = 0
	}
	_, ok = torrentMap["event"]
	if !ok {
		torrentMap["event"] = STOPPED
	}

	x := TorrentRequestData{
		torrentMap["info_hash"].(string),
		torrentMap["peer_id"].(string),
		torrentMap["ip"].(string),
		torrentMap["port"].(int),
		torrentMap["uploaded"].(int),
		torrentMap["downloaded"].(int),
		torrentMap["left"].(int),
		torrentMap["event"].(int),
	}
	return &x
}

func worker(torrdata *TorrentRequestData) {
	fmt.Println("%v", torrdata)
}

func requestHandler(w http.ResponseWriter, req *http.Request) {
	torrentdata := parseTorrentGetRequestURI(req.RequestURI)
	data := fillEmptyMapValues(torrentdata)
	go worker(data)
}

func RunServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/announce", requestHandler)
	http.ListenAndServe(":3000", mux)
}