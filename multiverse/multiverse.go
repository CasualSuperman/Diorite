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
		List      []*Card
	}
	Pronunciations trie.Trie
	Modified       time.Time
}

// Initialize the phonetics map for a constructed Multiverse.
func (m Multiverse) Initialize() {
	m.Pronunciations = generatePhoneticsMaps(m.Cards.List)
}
