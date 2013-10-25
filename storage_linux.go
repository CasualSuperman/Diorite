package main

import (
	"os"
)

var StorageDir string
var MultiverseFileName string

func init() {
	StorageDir = os.ExpandEnv("$HOME/.diorite/")
	MultiverseFileName = StorageDir + "multiverse.mtg"
}
