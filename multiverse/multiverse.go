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
		List      []scrubbedCard
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

type scrubbedCard struct {
	Card  *Card
	Ascii string
}

func scrubCards(list []*Card) []scrubbedCard {
	l := make([]scrubbedCard, len(list))

	for i, card := range list {
		l[i] = scrubbedCard{
			card,
			preventUnicode(card.Name),
		}
	}

	return l
}
