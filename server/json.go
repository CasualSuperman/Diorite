package main

import (
	"time"

	m "github.com/CasualSuperman/Diorite/multiverse"
	"github.com/glenn-brown/skiplist"
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
		if card.Card.Name == cardName {
			return i
		}
	}
	return -1
}

// Convert to a Multiverse.
func (om onlineMultiverse) Convert() (mv m.Multiverse) {
	mv.Sets = make([]*m.Set, 0, len(om.Sets))
	mv.Cards.Printings = skiplist.New()
	mv.Modified = om.Modified

	for _, set := range om.Sets {
		mv.Sets = append(mv.Sets, setFromJSON(set))
		for _, card := range set.Cards {
			var printing = struct {
				ID  m.MultiverseID
				Set *m.Set
			}{m.MultiverseID(card.MultiverseID), mv.Sets[len(mv.Sets)-1]}

			index := getCardIndex(mv.Cards.List, card.Name)
			if index == -1 {
				index = len(mv.Cards.List)
				c := new(m.Card)
				copyCardFields(&card, c)
				c.Printings = append(c.Printings, printing)
				mv.Cards.List.Add(c)
			} else {
				c := mv.Cards.List[index].Card
				c.Printings = append(c.Printings, printing)
			}
			mv.Cards.Printings.Insert(card.MultiverseID, index)
		}
	}

	return
}
