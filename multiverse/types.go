package multiverse

import (
	"time"

	"github.com/glenn-brown/skiplist"
)

// ManaColor is a bitmask of possible Mana Colors.
type ManaColor byte

// BorderColor indicates the color of card borders within a Set.
type BorderColor byte

// Rarity of a card for a specific printing.
type Rarity byte

// MultiverseID is a unique ID for a single printing of a card.
type MultiverseID int32

// SetType indicates the various set types.
type SetType byte

// The colors of mana that exist in the Multiverse.
var ManaColors = struct {
	Colorless, White, Blue, Black, Red, Green ManaColor
}{0, 1, 2, 4, 8, 16}

// The borders that cards have.
var BorderColors = struct {
	White, Black, Silver BorderColor
}{1, 2, 3}

// Rarities of cards.
var Rarities = struct {
	Common, Uncommon, Rare, Mythic, Basic, Special Rarity
}{1, 2, 3, 4, 5, 6}

// Set types.
var SetTypes = struct {
	Core, Expansion, Reprint, Box, Un, FromTheVault, PremiumDeck, DuelDeck, Starter, Commander, Planechase, Archenemy SetType
}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

// Multiverse is an entire Magic: The Gathering multiverse.
// It contains the available cards, sets, formats, and legality information, as well as ways to interpret, manipulate, and filter that data.
type Multiverse struct {
	Sets  []*Set
	Cards struct {
		Printings *skiplist.T
		List      CardList
	}
	Modified time.Time
}

// Set is a Magic: The Gathering set, such as Innistrad or Zendikar.
type Set struct {
	Name     string
	Code     string
	Released time.Time
	Border   BorderColor
	Type     SetType
	Block    string
	Cards    []MultiverseID
}

// Card is a Magic: The Gathering card, such as Ætherling or Blightning.
type Card struct {
	Name   string
	Cmc    float32
	Cost   string
	Colors ManaColor

	Supertypes, Types []string

	Text   string
	Flavor string

	Artist string
	Number string

	Power, Toughness struct {
		Val      float32
		Original string
	}

	Rulings []Ruling

	Printings []Printing

	Restricted, Banned []*Format
}

type Printing struct {
	ID     MultiverseID
	Set    *Set
	Rarity Rarity
}

// Ruling is a ruling made by a judge that can clarify difficult situations that may arise.
type Ruling struct {
	Date time.Time
	Text string
}

// IsCreature is a convenience method that returns if the card is a creature.
func (c *Card) IsCreature() bool {
	for _, supertype := range c.Supertypes {
		if supertype == "Creature" {
			return true
		}
	}
	return false
}
