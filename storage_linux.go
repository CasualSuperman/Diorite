package main

import (
	"os"
	"strings"
)

func init() {
	StorageDir = strings.Join([]string{os.ExpandEnv("$HOME"), ".diorite"}, string(os.PathSeparator))
	MultiverseFileName = StorageDir + string(os.PathSeparator) + "multiverse.mtg"
}
