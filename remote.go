package main

import (
	"encoding/json"
	"net/http"
	"time"
)

const remoteDBLocation = "http://mtgjson.com/json/AllSets-x.json"
const lastModifiedFormat = time.RFC1123

func onlineModifiedAt() (time.Time, error) {
	var netClient http.Client

	resp, err := netClient.Head(remoteDBLocation)
	if err != nil {
		return time.Time{}, err
	}
	remoteModified := resp.Header.Get("Last-Modified")
	rModTime, _ := time.Parse(lastModifiedFormat, remoteModified)
	if err != nil {
		return time.Time{}, err
	}
	return rModTime, nil
}

func downloadOnline() (map[string]jsonSet, error) {
	resp, err := http.Get(remoteDBLocation)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	structure := make(map[string]jsonSet)

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&structure)

	return structure, err
}
