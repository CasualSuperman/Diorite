package main

import (
	"os"
)

func init() {
	StorageDir = os.ExpandEnv("${APPDATA}/CasualSuperman/Diorite/")
	MultiverseFileName = StorageDir + "multiverse.mtg"
}
