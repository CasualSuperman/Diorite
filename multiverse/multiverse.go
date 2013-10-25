package multiverse

import (
	"time"

	"github.com/CasualSuperman/Diorite/trie"
	"github.com/glenn-brown/skiplist"
)

// Multiverse is an entire Magic: The Gathering multiverse.
// It contains the available cards, sets, formats, and legality information, as well as ways to interpret, manipulate, and filter that data.
type Multiverse struct {
	Sets           map[string]*Set
	Cards          *skiplist.T
	cardList       []*Card
	Pronunciations *trie.Trie
	Modified       time.Time
}

func getCardIndex(cardList []*Card, cardName string) int {
	for i, card := range cardList {
		if card.Name == cardName {
			return i
		}
	}
	return -1
}

// Create a Multiverse from the provided map of sets, with a modification time of the provided time.
func Create(json map[string]jsonSet, modified time.Time) (m Multiverse) {
	m.Sets = make(map[string]*Set)
	m.Cards = skiplist.New()
	m.cardList = make([]*Card, 0)
	m.Modified = modified

	for _, set := range json {
		m.Sets[set.Name] = setFromJSON(set)
		for _, card := range set.Cards {
			index := getCardIndex(m.cardList, card.Name)
			if index == -1 {
				index = len(m.cardList)
				c := new(Card)
				copyCardFields(&card, c)
				m.cardList = append(m.cardList, c)
			}
			m.Cards.Insert(card.MultiverseID, index)
		}
	}

	m.Pronunciations = generatePhoneticsMaps(m.cardList)

	return
}
