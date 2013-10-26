package main

import (
	"os"
)

func init() {
	StorageDir = os.ExpandEnv("$HOME/Library/Application Support/Diorite/")
	MultiverseFileName = StorageDir + "multiverse.mtg"
}
