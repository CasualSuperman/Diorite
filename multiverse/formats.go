package multiverse

import (
	"time"
)

var Formats = []struct {
	Name string
}{}

type legalityCardCheck func(*Card) bool
type legalitySetCheck func(*Set) bool

func unSet(s *Set) bool {
	return s.Type == SetType.Un
}

func vintageSetLegal(s *Set) bool {
	return !unSet(s)
}

func legacySetLegal(s *Set) bool {
	return !unSet(s)
}

var firstModernSet = time.Date(2003, time.July, 28, 0, 0, 0, 0, time.UTC)

func modernSetLegal(s *Set) bool {
	return !unSet(s) && !s.Released.Before(firstModernSet)
}
