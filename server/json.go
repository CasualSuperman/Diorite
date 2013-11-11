package main

import (
	"sort"
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

func getCardIndex(cardList m.CardList, cardName string) int {
	for i, card := range cardList {
		if card.Name == cardName {
			return i
		}
	}
	return -1
}

type setSorter struct {
	sets []*m.Set
	by   func(s1, s2 *m.Set) bool
}

func (s setSorter) Len() int {
	return len(s.sets)
}

func (s setSorter) Swap(i, j int) {
	s.sets[i], s.sets[j] = s.sets[j], s.sets[i]
}

func (s setSorter) Less(i, j int) bool {
	return s.by(s.sets[i], s.sets[j])
}

// Convert to a Multiverse.
func (om onlineMultiverse) Convert() (mv m.Multiverse) {
	mv.Sets = make([]*m.Set, 0, len(om.Sets))
	mv.Modified = om.Modified

	for _, set := range om.Sets {
		mSet := setFromJSON(set)
		mv.Sets = append(mv.Sets, mSet)
	}

	s := setSorter{
		mv.Sets,
		releaseDateSort,
	}

	sort.Sort(s)

	for _, set := range om.Sets {
		var setIndex int

		for i, s := range mv.Sets {
			if s.Name == set.Name {
				setIndex = i
				break
			}
		}

		for _, card := range set.Cards {
			var rarity m.Rarity
			switch card.Rarity {
			case "Common":
				rarity = m.Rarities.Common
			case "Uncommon":
				rarity = m.Rarities.Uncommon
			case "Rare":
				rarity = m.Rarities.Rare
			case "Mythic Rare":
				rarity = m.Rarities.Mythic
			case "Special":
				rarity = m.Rarities.Special
			case "Basic Land":
				rarity = m.Rarities.Basic
			}

			var printing = m.Printing{
				m.MultiverseID(card.MultiverseID),
				mv.Sets[setIndex],
				rarity,
			}

			index := getCardIndex(mv.Cards, card.Name)
			if index == -1 {
				index = len(mv.Cards)
				c := new(m.Card)
				copyCardFields(&card, c)
				c.Printings = append(c.Printings, printing)
				mv.Cards.Add(*c)
			} else {
				c := &mv.Cards[index]
				c.Printings = append(c.Printings, printing)
			}
		}
	}

	return
}

func releaseDateSort(s1, s2 *m.Set) bool {
	return s2.Released.Before(s1.Released)
}
