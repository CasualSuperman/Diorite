package multiverse

import (
	"time"

	"github.com/CasualSuperman/Diorite/trie"
	"github.com/glenn-brown/skiplist"
)

// Multiverse is an entire Magic: The Gathering multiverse.
// It contains the available cards, sets, formats, and legality information, as well as ways to interpret, manipulate, and filter that data.
type Multiverse struct {
	Sets  map[string]*Set
	Cards struct {
		Printings *skiplist.T
		List      CardList
	}
	Pronunciations trie.Trie
	Modified       time.Time
}

// Initialize the phonetics map for a constructed Multiverse.
func (m Multiverse) Initialize() {
	m.Pronunciations = generatePhoneticsMaps(m.Cards.List)
}

func (m Multiverse) Card(id int) *Card {
	return m.Cards.List[id].Card
}

type CardList []scrubbedCard

func (c *CardList) Add(candidate *Card) int {
	for i, card := range *c {
		if card.Card == candidate {
			return i
		}
	}

	*c = append(*c, scrubbedCard{
		candidate,
		preventUnicode(candidate.Name),
	})

	return len(*c) - 1
}

func (c *CardList) Len() int {
	return len(*c)
}

type scrubbedCard struct {
	Card  *Card
	Ascii string
}

func scrubCards(list []*Card) CardList {
	l := CardList(make([]scrubbedCard, len(list)))

	for i, card := range list {
		l[i] = scrubbedCard{
			card,
			preventUnicode(card.Name),
		}
	}

	return l
}
