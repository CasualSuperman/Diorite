package main

import (
	"os"
	"strings"
)

func init() {
	StorageDir = strings.Join([]string{os.ExpandEnv("${APPDATA}"), "CasualSuperman", "Diorite"}, string(os.PathSeparator))
	MultiverseFileName = StorageDir + string(os.PathSeparator) + "multiverse.mtg"
}
