package main

import (
	"os"
)

var StorageDir string
var MultiverseFileName string

func init() {
	StorageDir = os.ExpandEnv("$HOME/Library/Application Support/Diorite/")
	MultiverseFileName = StorageDir + "multiverse.mtg"
}
