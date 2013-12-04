package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"runtime"
	"time"

	m "github.com/CasualSuperman/Diorite/multiverse"
)

var local = flag.Bool("local", false, "Connect to a server running on localhost.")
var debug = flag.Bool("debug", false, "Keep the server alive even if the page unloads.")

type exitSignal struct {
	code   int
	reason string
}

func exit(es exitSignal) {
	log.Println(es.reason)
	os.Exit(es.code)
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU() * 8)

	multiverseChan := make(chan *m.Multiverse)
	exitChan := make(chan exitSignal)

	multiverseDirAvailable := makeStorageDir()

	var multiverse *m.Multiverse

	if multiverseDirAvailable {
		multiverse = loadMultiverse()
	}

	server := NewServer(multiverse)

	if !*local {
		log.Println("Starting display server.")
		go server.Serve(":7000", exitChan)
	}

	log.Println("Checking for updates.")
	go watchForUpdates(multiverse, multiverseChan)

	for {
		select {
		case multiverse = <-multiverseChan:
			log.Println("Multiverse updated.")
			if !*local {
				server.UpdateMultiverse(multiverse)
			} else {
				go func() {
					exitChan <- exitSignal{0, "Successfully updated."}
				}()
			}
		case code := <-exitChan:
			exit(code)
		}
	}
}

func watchForUpdates(currMultiverse *m.Multiverse, ch chan *m.Multiverse) {
	updateFrequency := time.Tick(5 * time.Minute)
	for {
		newMultiverse := downloadUpdates(currMultiverse)
		if newMultiverse.Modified.After(currMultiverse.Modified) {
			currMultiverse = newMultiverse
			ch <- currMultiverse
		}
		<-updateFrequency
	}
}

func loadLocalMultiverse(location string) (m.Multiverse, error) {
	multiverseFile, err := os.Open(location)

	if err != nil {
		return m.Multiverse{}, err
	}

	defer multiverseFile.Close()

	return m.Read(multiverseFile)
}

func connectToServer() (serverConnection, error) {
	if *local {
		return connectToLocalServer()
	} else {
		return connectToDefaultServer()
	}
}

func makeStorageDir() bool {
	err := os.MkdirAll(StorageDir, os.ModePerm|os.ModeDir)

	if err != nil {
		log.Println("Unable to create application storage directory. Multiverse will not be saved.")
	}

	return err == nil
}

func loadMultiverse() *m.Multiverse {
	log.Println("Loading local multiverse.")

	multiverse, err := loadLocalMultiverse(MultiverseFileName)

	if err != nil {
		if os.IsNotExist(err) || os.IsPermission(err) {
			log.Println("No local database available. A local copy will be downloaded.")
		} else {
			log.Printf("Unable to load multiverse: %s\n", err)
		}
	} else {
		log.Println("Multiverse loaded.")
	}

	return &multiverse
}

func downloadUpdates(multiverse *m.Multiverse) *m.Multiverse {
	log.Println("Contacting update server.")

	server, err := connectToServer()

	if err != nil {
		if !multiverse.Loaded() {
			log.Println("No local multiverse available, and unable to download copy. Unable to continue.")
		} else {
			log.Println("Warning: Online database unavailable. Your card index may be out of date.")
		}
		return multiverse
	}

	defer server.Close()

	if multiverse.Loaded() {
		log.Println("Checking for multiverse updates.")
	} else {
		log.Println("Downloading online multiverse.")
	}

	if multiverse.Loaded() && !server.Modified().After(multiverse.Modified) {
		log.Println("No updates found.")
		return multiverse
	}

	log.Println("Updates available. Downloading.")

	data := server.RawMultiverse()
	buf := bytes.NewBuffer(data)

	saveTo, err := os.Create(MultiverseFileName)

	if err != nil {
		log.Println("Unable to save update to multiverse. Continuing, but it will be redownloaded on next startup.")
	} else {
		_, err := saveTo.Write(data)
		if err != nil && multiverse.Loaded() {
			log.Println("Error saving multiverse. Rolling back changes.")
			saveTo.Truncate(0)
			multiverse.Write(saveTo)
			saveTo.Close()
		} else if err != nil {
			log.Println("Error saving Multiverse. Removing.")
			saveTo.Close()
			os.Remove(MultiverseFileName)
		}
	}

	newM, err := m.Read(buf)

	if err != nil {
		log.Printf("Error updating: %s\n", err)
		if !multiverse.Loaded() {
			log.Println("Downloading multiverse failed and no local database available. Unable to continue.")
		} else {
			log.Println("Unable to download most recent multiverse. Continuing with an out-of-date version.")
		}
	} else {
		log.Println("Multiverse downloaded!")
		return &newM
	}
	return multiverse
}
