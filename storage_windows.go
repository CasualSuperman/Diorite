package main

import (
	"os"
)

var StorageDir string
var MultiverseFileName string

func init() {
	StorageDir = os.ExpandEnv("${APPDATA}/CasualSuperman/Diorite/")
	MultiverseFileName = StorageDir + "multiverse.mtg"
}
