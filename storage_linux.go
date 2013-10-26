package main

import (
	"os"
)

func init() {
	StorageDir = os.ExpandEnv("$HOME/.diorite/")
	MultiverseFileName = StorageDir + "multiverse.mtg"
}
