package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	m "github.com/CasualSuperman/Diorite/multiverse"
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

	for _, format := range m.Formats.List {
		log.Println("Downloading " + format.Name + " banned list.")
		if banned, err := getBannedList(format.Name); err == nil {
		bannedLoop:
			for _, name := range banned {
				for i := range multiverse.Cards.List {
					var card *m.Card = multiverse.Cards.List[i].Card
					if card.Name == name {
						card.Banned = append(card.Banned, format)
						continue bannedLoop
					}
				}
			}
		}

		log.Println("Downloading " + format.Name + " restricted list.")
		if restricted, err := getRestrictedList(format.Name); err == nil {
		restrictedLoop:
			for _, name := range restricted {
				for i := range multiverse.Cards.List {
					var card *m.Card = multiverse.Cards.List[i].Card
					if card.Name == name {
						card.Restricted = append(card.Restricted, format)
						continue restrictedLoop
					}
				}
			}
		}
	}

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
