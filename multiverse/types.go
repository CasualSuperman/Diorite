package multiverse

import (
	"time"
)

type ManaColor byte
type BorderColor byte
type Rarity byte
type MultiverseID int32
type SetType byte

// The colors of mana that exist in the Multiverse.
var ManaColors = struct {
	White, Blue, Black, Red, Green ManaColor
}{1, 2, 4, 8, 16}

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

	Rarity Rarity

	Text   string
	Flavor string

	Artist string
	Number string

	Power, Toughness struct {
		Val      float32
		Original string
	}

	Rulings []Ruling
}

type Ruling struct {
	Date time.Time
	Text string
}

type setSorter struct {
	sets []Set
	by   func(s1, s2 *Set) bool
}

func (s *setSorter) Len() int {
	return len(s.sets)
}

func (s *setSorter) Swap(i, j int) {
	s.sets[i], s.sets[j] = s.sets[j], s.sets[i]
}

func (s *setSorter) Less(i, j int) bool {
	return s.by(&s.sets[i], &s.sets[j])
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
