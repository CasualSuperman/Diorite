package main

import (
	"os"
	"path/filepath"
)

var defaultLocation string

func init() {
	if os.ExpandEnv("${XDG_DATA_HOME}") != "" {
		defaultLocation = os.ExpandEnv("${XDG_DATA_HOME}/diorite/multiverse.mtg")
	} else {
		defaultLocation = os.ExpandEnv("${HOME}/.local/share/diorite/multiverse.mtg")
	}
	MultiverseFileName = filepath.FromSlash(defaultLocation)
	StorageDir = filepath.Dir(MultiverseFileName)
}
