package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	m "github.com/CasualSuperman/Diorite/multiverse"
)

var local = flag.Bool("local", false, "Connect to a server running on localhost.")

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU() * 8)

	var multiverse m.Multiverse

	multiverseDirAvailable := makeStorageDir()

	if multiverseDirAvailable {
		multiverse = loadMultiverse()
	}

	multiverse = downloadUpdates(multiverse)

	if !multiverse.Loaded() {
		os.Exit(1)
	}

	log.Println("Cards in multiverse:", len(multiverse.Cards))

	searchTerm := ""
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i][0] != '-' {
			if len(searchTerm) > 0 {
				searchTerm += " "
			}
			searchTerm += os.Args[i]
		}
	}

	if len(searchTerm) > 0 {
		card := multiverse.FuzzyNameSearch(searchTerm, 1)
		if len(card) > 0 {
			fmt.Println(card[0])
		}
	}
}

func loadLocalMultiverse(location string) (m.Multiverse, error) {
	multiverseFile, err := os.Open(location)
	defer multiverseFile.Close()

	if err != nil {
		return m.Multiverse{}, err
	}

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

func loadMultiverse() m.Multiverse {
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

	return multiverse
}

func downloadUpdates(multiverse m.Multiverse) m.Multiverse {
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
		return newM
	}
	return multiverse
}
