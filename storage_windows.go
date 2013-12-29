package main

import (
	"os"
	"path/filepath"
)

var defaultLocation = os.ExpandEnv("${APPDATA}/CasualSuperman/Diorite/multiverse.mtg")

func init() {
	MultiverseFileName = filepath.FromSlash(defaultLocation)
	StorageDir = filepath.Dir(MultiverseFileName)
}
