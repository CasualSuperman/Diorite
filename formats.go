package main

import (
	"time"
)

var Formats = []struct {
	Name string
}{}

type legalityCardCheck func(*card) bool
type legalitySetCheck func(*set) bool

func unSet(s *set) bool {
	return s.Type == SetType.Un
}

func vintageSetLegal(s *set) bool {
	return !unSet(s)
}

func legacySetLegal(s *set) bool {
	return !unSet(s)
}

var firstModernSet = time.Date(2003, time.July, 28, 0, 0, 0, 0, time.UTC)

func modernSetLegal(s *set) bool {
	return !unSet(s) && !s.Released.Before(firstModernSet)
}
