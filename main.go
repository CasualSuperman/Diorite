package main

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"runtime/pprof"
	"time"

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

		f, _ := os.Create("inflateprofile")
		pprof.StartCPUProfile(f)
		m = multiverse.Inflate(multiverseFile)
		pprof.StopCPUProfile()
		fmt.Println("Multiverse loaded.")
		multiverseLoaded = true
	}

	fmt.Println("Checking for multiverse updates.")
	mostRecentUpdate, err := onlineModifiedAt()

	if err != nil {
		fmt.Println("Warning! Online database unavailable. Your card index may be out of date.")
	}

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
		file, err := os.Create(multiverseFileName)
		if err == nil {
			defer file.Close()
			fmt.Println("Saving downloaded multiverse.")
			err := om.WriteTo(file)
			if err != nil {
				fmt.Println("Error saving multiverse:", err)
			}
		} else {
			fmt.Println(err)
		}
		fmt.Println("Transforming multiverse.")
		m = multiverse.Create(om.Sets, time.Now())
	} else {
		fmt.Println("No updates available.")
	}

	f, _ := os.Create("searchprofile")
	pprof.StartCPUProfile(f)
	m.SearchByName("a")
	m.SearchByName("av")
	m.SearchByName("ava")
	m.SearchByName("avat")
	m.SearchByName("avata")
	m.SearchByName("avatar")
	cards := m.SearchByName("lightning")
	pprof.StopCPUProfile()
	names := make([]string, len(cards))
	for i, card := range cards {
		names[i] = card.Name
	}
	fmt.Println(names)

	f, _ = os.Create("memprofile")
	pprof.WriteHeapProfile(f)
	f.Close()
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
