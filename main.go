package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"time"
)

const dataLocation = ".diorite"
const dbName = "cards.json"
const remoteDBLocation = "http://mtgjson.com/json/AllSets-x.json"
const lastModifiedFormat = time.RFC1123

var netClient http.Client

func main() {
	// Get the information about our user.
	u, err := user.Current()
	if err != nil {
		// Well that ended poorly.
		fmt.Println("Something went horribly wrong.")
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

	if staleDB() {
		fmt.Println("Updating card database.")
		err := updateDB()
		if err != nil {
			fmt.Println("Unable to update database.")
		}
	} else {
		fmt.Println("Card database is up-to-date.")
	}
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

func staleDB() bool {
	resp, err := netClient.Head(remoteDBLocation)
	remoteModified := resp.Header.Get("Last-Modified")
	rModTime, err := time.Parse(lastModifiedFormat, remoteModified)
	if err != nil {
		return true
	}
	localDB, err := os.Stat(dbName)
	if err != nil {
		return true
	}
	return rModTime.After(localDB.ModTime())
}

func updateDB() error {
	resp, err := netClient.Get(remoteDBLocation)
	if err != nil {
		return err
	}

	file, err := os.Create(dbName)
	if err != nil {
		return err
	}

	io.Copy(file, resp.Body)
	resp.Body.Close()
	file.Close()
	return nil
}
