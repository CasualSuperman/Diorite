package main

import (
	"fmt"
	"math"
	"os"
	"os/user"
	"runtime"
	"sync"

	"github.com/CasualSuperman/Diorite/multiverse"
)

const dataLocation = ".diorite"
const multiverseFileName = "multiverse.mtg"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Get the information about our user.
	u, err := user.Current()
	if err != nil {
		// Well that ended poorly.
		fmt.Println("Something went horribly wrong.")
		os.Exit(1)
	}

	// Find out if our storage directory exists.
	if !homeDirExists(u.HomeDir) {
		fmt.Println("Creating ~/" + dataLocation)
		err := createHomeDir(u.HomeDir)

		if err != nil {
			fmt.Println("Unable to create required storage directory.")
		}
	}

	os.Chdir(dataDir(u.HomeDir))
	multiverseFile, err := os.Open(multiverseFileName)

	var m multiverse.Multiverse
	multiverseLoaded := false

	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			fmt.Println("No local database available. A local copy will be downloaded.")
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Loading local multiverse.")
		m, err = multiverse.Read(multiverseFile)
		multiverseFile.Close()
	}

	if err != nil {
		fmt.Println("Unable to load multiverse:", err)
	} else {
		fmt.Println("Multiverse loaded.")
		multiverseLoaded = true
	}

	fmt.Println("Checking for multiverse updates.")
	mostRecentUpdate, err := onlineModifiedAt()

	if err != nil {
		fmt.Println("Warning! Online database unavailable. Your card index may be out of date.")
	}

	var saved sync.WaitGroup

	if mostRecentUpdate.After(m.Modified) {
		fmt.Println("Multiverse update available! Downloading now.")
		om, err := downloadOnline()
		if err != nil {
			if !multiverseLoaded {
				fmt.Println("Unable to download multiverse and no local database available. Unable to continue.")
				os.Exit(1)
			}
			fmt.Println("Unable to download most recent multiverse. Continuing with an out-of-date version.")
		}
		fmt.Println("Transforming multiverse.")
		m = multiverse.Create(om.Sets, om.Modified)

		file, err := os.Create(multiverseFileName)

		if err != nil {
			fmt.Println(err)
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

func dataDir(baseDir string) string {
	return baseDir + string(os.PathSeparator) + dataLocation
}

func homeDirExists(baseDir string) bool {
	dir, err := os.Stat(dataDir(baseDir))
	return err == nil && dir.IsDir()
}

func createHomeDir(baseDir string) error {
	return os.Mkdir(dataDir(baseDir), os.ModePerm|os.ModeDir)
}
