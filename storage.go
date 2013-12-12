package main

import "path/filepath"

// The location for storing the multiverse.
var MultiverseFileName = filepath.FromSlash(defaultLocation)
// The folder that should contain the multiverse.
var StorageDir = filepath.Dir(MultiverseFileName)
