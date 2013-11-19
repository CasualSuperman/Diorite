package main

import (
	"os"
	"strings"
)

func init() {
	StorageDir = strings.Join([]string{os.ExpandEnv("$HOME"), "Library", "Application Support", "Diorite"}, string(os.PathSeparator))
	MultiverseFileName = StorageDir + string(os.PathSeparator) + "multiverse.mtg"
}
