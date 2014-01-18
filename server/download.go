package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
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

func getFormatLists(ret chan formatList, errChan chan error) {
	log.Println("Downloading banned/restricted lists.")
	var externalGroup sync.WaitGroup

	for _, format := range m.Formats.List {
		externalGroup.Add(1)
		go func(format *m.Format) {
			defer externalGroup.Done()
			list := formatList{
				format,
				make(map[string]bool),
				make(map[string]bool),
			}

			var internalGroup sync.WaitGroup
			internalGroup.Add(1)

			go func() {
				defer internalGroup.Done()
				if banned, err := getBanList(format.Name); err == nil {
					for _, name := range banned {
						list.Banned[name] = true
					}
				} else {
					errChan <- err
				}
			}()

			if format == m.Formats.Vintage {
				internalGroup.Add(1)
				go func() {
					defer internalGroup.Done()
					if restricted, err := getRestrictList(format.Name); err == nil {
						for _, name := range restricted {
							list.Restricted[name] = true
						}
					} else {
						errChan <- err
					}
				}()
			}

			internalGroup.Wait()
			ret <- list
		}(format)
	}

	externalGroup.Wait()

	log.Println("Banned/restricted lists finished downloading.")

	close(ret)
}

func getMultiverse(ret chan m.Multiverse, errChan chan error) {
	var structure onlineMultiverse

	resp, err := getOnline("GET")

	log.Println("Multiverse downloaded. Converting JSON.")

	if err != nil {
		errChan <- err
		return
	}

	dec := json.NewDecoder(resp.Body)

	if err = dec.Decode(&structure.Sets); err != nil {
		errChan <- err
		return
	}

	resp.Body.Close()

	log.Println("JSON converted.")

	remoteModified := resp.Header.Get("Last-Modified")
	rModTime, err := time.Parse(lastModifiedFormat, remoteModified)

	if err != nil {
		errChan <- err
		return
	}

	structure.Modified = rModTime

	log.Println("Converting JSON structure to Multiverse.")

	multiverse := structure.Convert()

	log.Println("Structure converted.")

	ret <- multiverse

	structure.Reset()
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
