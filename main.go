package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"

	m "github.com/CasualSuperman/Diorite/multiverse"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	canSaveMultiverse := true
	multiverseLoaded := false

	err := os.MkdirAll(StorageDir, os.ModePerm|os.ModeDir)

	if err != nil {
		log.Println("Unable to create application storage directory. Multiverse will not be saved.")
		canSaveMultiverse = false
	}

	var multiverse m.Multiverse

	if canSaveMultiverse {
		multiverseFile, err := os.Open(MultiverseFileName)

		if err != nil {
			if os.IsNotExist(err) || os.IsPermission(err) {
				log.Println("No local database available. A local copy will be downloaded.")
			} else {
				log.Fatalln(err)
			}
		} else {
			log.Println("Loading local multiverse.")
			multiverse, err = m.Read(multiverseFile)
			multiverseFile.Close()

			if err != nil {
				log.Printf("Unable to load multiverse: %s\n", err)
			} else {
				log.Println("Multiverse loaded.")
				multiverseLoaded = true
			}
		}

	}

	if multiverseLoaded {
		log.Println("Checking for multiverse updates.")
	} else {
		log.Println("Downloading online multiverse.")
	}
	mostRecentUpdate, err := onlineModifiedAt()

	if err != nil {
		log.Println(err)
		if multiverseLoaded == false {
			log.Fatalln("No local multiverse available, and unable to download copy. Unable to continue.")
		}
		log.Println("Warning: Online database unavailable. Your card index may be out of date.")
	}

	if mostRecentUpdate.After(multiverse.Modified) {
		var saveTo io.Writer
		if canSaveMultiverse {
			saveTo, err = os.Create(MultiverseFileName)
			if err != nil {
				log.Println("Unable to save update to multiverse. Continuing, but it will be redownloaded on next startup.")
				saveTo = nil
			}
		}
		log.Println("Multiverse update available! Downloading now.")
		newM, err := downloadMultiverse(saveTo)
		if err != nil {
			log.Printf("Error downloading: %s\n", err)
			if !multiverseLoaded {
				log.Fatalln("Downloading multiverse failed and no local database available. Unable to continue.")
			}
			log.Println("Unable to download most recent multiverse. Continuing with an out-of-date version.")
		} else {
			multiverse = newM
		}
	} else {
		log.Println("No updates available.")
	}

	cards := multiverse.FuzzyNameSearch("aetherling", 15)

	names := make([]string, len(cards))
	for i, card := range cards {
		names[i] = card.Name
		fmt.Println(card.Name)

		if math.IsNaN(float64(card.Toughness.Val)) {
			if card.IsCreature() {
				fmt.Println(card.Toughness.Original)
			}
		} else {
			fmt.Println(card.Toughness.Val)
		}
	}
	fmt.Println(names)
}
