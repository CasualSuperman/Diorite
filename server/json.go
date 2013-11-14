package main

import (
	"time"

	m "github.com/CasualSuperman/Diorite/multiverse"
)

type jsonRuling struct {
	Date string `json:"date"`
	Text string `json:"text"`
}

type jsonCard struct {
	Name string `json:"name"`

	ManaCost string   `json:"manaCost"`
	Cmc      float32  `json:"cmc"`
	Colors   []string `json:"colors"`

	Type       string   `json:"type"`
	Supertypes []string `json:"supertypes"`
	Types      []string `json:"types"`
	Subtypes   []string `json:"subtypes"`

	Rarity string `json:"rarity"`

	Text string `json:"text"`

	Flavor string `json:"flavor"`

	Artist string `json:"artist"`
	Number string `json:"number"`

	Power     string `json:"power"`
	Toughness string `json:"toughness"`

	Layout       string `json:"layout"`
	MultiverseID int    `json:"multiverseid"`

	ImageName string `json:"imageName"`

	Rulings []jsonRuling `json:"rulings"`
}

type jsonSet struct {
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	ReleaseDate string     `json:"releaseDate"`
	Border      string     `json:"border"`
	Type        string     `json:"type"`
	Block       string     `json:"block"`
	Cards       []jsonCard `json:"cards"`
}

// OnlineMultiverse is a convenience type for the conversion type from mtgjson.com.
// Eventually this will be removed.
type onlineMultiverse struct {
	Sets     map[string]jsonSet
	Modified time.Time
}

type setSorter struct {
	sets []m.Set
	by   func(s1, s2 *m.Set) bool
}

func (s setSorter) Len() int {
	return len(s.sets)
}

func (s setSorter) Swap(i, j int) {
	s.sets[i], s.sets[j] = s.sets[j], s.sets[i]
}

func (s setSorter) Less(i, j int) bool {
	return s.by(&s.sets[i], &s.sets[j])
}

func releaseDateSort(s1, s2 *m.Set) bool {
	return s2.Released.Before(s1.Released)
}
