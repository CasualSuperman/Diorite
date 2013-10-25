package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/CasualSuperman/Diorite/multiverse"
)

const remoteDBLocation = "http://mtgjson.com/json/AllSets-x.json"
const lastModifiedFormat = time.RFC1123

func onlineModifiedAt() (time.Time, error) {
	resp, err := getOnline("HEAD")

	if err != nil {
		return time.Time{}, err
	}
	remoteModified := resp.Header.Get("Last-Modified")
	return time.Parse(lastModifiedFormat, remoteModified)
}

func downloadMultiverse() (m multiverse.Multiverse, err error) {
	var structure multiverse.OnlineMultiverse
	resp, err := getOnline("GET")

	if err != nil {
		return
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	if err = dec.Decode(&structure.Sets); err != nil {
		return
	}

	remoteModified := resp.Header.Get("Last-Modified")
	rModTime, err := time.Parse(lastModifiedFormat, remoteModified)
	if err != nil {
		return
	}
	structure.Modified = rModTime

	return structure.Convert(), err
}

func getOnline(method string) (*http.Response, error) {
	var netClient http.Client

	switch method {
	case "HEAD":
		return netClient.Head(remoteDBLocation)
	case "GET":
		return netClient.Get(remoteDBLocation)
	default:
		return nil, errors.New("unknown method")
	}
}
