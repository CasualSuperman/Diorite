package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"sync"

	"github.com/CasualSuperman/Diorite/multiverse"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	canSaveMultiverse := true
	multiverseLoaded := false

	err := os.MkdirAll(StorageDir, os.ModePerm|os.ModeDir)

	if err != nil {
		fmt.Errorf("Unable to create application storage directory. Multiverse will not be saved.")
		canSaveMultiverse = false
	}

	var m multiverse.Multiverse

	if canSaveMultiverse {
		multiverseFile, err := os.Open(MultiverseFileName)

		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No local database available. A local copy will be downloaded.")
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Loading local multiverse.")
			m, err = multiverse.Read(multiverseFile)
			multiverseFile.Close()

			if err != nil {
				fmt.Errorf("Unable to load multiverse: %s", err)
			} else {
				fmt.Println("Multiverse loaded.")
				multiverseLoaded = true
			}
		}

	}

	fmt.Println("Checking for multiverse updates.")
	mostRecentUpdate, err := onlineModifiedAt()

	if err != nil {
		if multiverseLoaded == false {
			fmt.Errorf("No local multiverse available, and unable to download copy. Unable to continue.")
			os.Exit(1)
		}
		fmt.Errorf("Warning: Online database unavailable. Your card index may be out of date.")
	}

	var saved sync.WaitGroup

	if mostRecentUpdate.After(m.Modified) {
		fmt.Println("Multiverse update available! Downloading now.")
		newM, err := downloadMultiverse()
		if err != nil {
			if !multiverseLoaded {
				fmt.Errorf("Downloading multiverse failed and no local database available. Unable to continue.")
				os.Exit(1)
			}
			fmt.Println("Unable to download most recent multiverse. Continuing with an out-of-date version.")
		} else {
			m = newM
		}

		if canSaveMultiverse {
			file, err := os.Create(MultiverseFileName)

			if err != nil {
				fmt.Errorf("Unable to save update to multiverse. Continuing with up-to-date multiverse, but it will be redownloaded on next startup.")
				fmt.Errorf("(Reason for failure: %s)", err)
			} else {
				saved.Add(1)
				go func() {
					defer file.Close()
					defer saved.Done()
					fmt.Println("Saving downloaded multiverse.")
					err := m.Write(file)
					if err != nil {
						fmt.Println("Error saving multiverse:", err)
					}
				}()
			}
		}
	} else {
		fmt.Println("No updates available.")
	}

	cards := m.FuzzyNameSearch("aetherling", 15)

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

	saved.Wait()
}
