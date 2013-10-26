package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
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

func getMultiverseData() ([]byte, time.Time, error) {
	var structure onlineMultiverse
	resp, err := getOnline("GET")

	if err != nil {
		return nil, time.Time{}, err
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	if err = dec.Decode(&structure.Sets); err != nil {
		return nil, time.Time{}, err
	}

	remoteModified := resp.Header.Get("Last-Modified")
	rModTime, err := time.Parse(lastModifiedFormat, remoteModified)
	if err != nil {
		return nil, time.Time{}, err
	}
	structure.Modified = rModTime

	multiverse := structure.Convert()

	var b bytes.Buffer
	multiverse.Write(&b)
	multiverseDL := b.Bytes()
	b.Reset()
	multiverseMod := multiverse.Modified

	return multiverseDL, multiverseMod, err
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
