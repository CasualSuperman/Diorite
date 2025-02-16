package multiverse

import (
	"time"
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

// SuperType indicates the various SuperTypes.
type SuperType byte

// Type indicates the various Types.
type Type byte

// Printing represents a specific printing of a card, Cancel from M10 is different from Cancel from M11.
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

// Multiverse is an entire Magic: The Gathering multiverse.
// It contains the available cards, sets, formats, and legality information, as well as ways to interpret, manipulate, and filter that data.
type Multiverse struct {
	Sets     []Set
	Cards    []Card
	Modified time.Time
}

// Set is a Magic: The Gathering set, such as Innistrad or Zendikar.
type Set struct {
	Name          string
	Code          string
	Released      time.Time
	Border        BorderColor
	Type          SetType
	Block         string
	Cards         []MultiverseID
	standardLegal bool
	extendedLegal bool
}

// Card is a Magic: The Gathering card, such as Ætherling or Blightning.
type Card struct {
	Name       string
	ascii      string
	metaphones []string
	Cmc        float32
	Cost       string
	Colors     ManaColor

	Supertypes SuperType
	Types      Type
	Subtypes   []string

	Text   string
	Flavor string

	Artist string
	Number string

	Power, Toughness usuallyNumeric
	Loyalty          int

	Rulings []Ruling

	Printings []Printing `json:"-"`

	Restricted, Banned []*Format
}

type usuallyNumeric struct {
	Val      float32
	Original string
}

func (c *Card) Is(f Filter) bool {
	ok, err := f.Ok(c)

	return ok && err == nil
}

func (s SetType) isTournamentLegal() bool {
	return s == SetTypes.Expansion || s == SetTypes.Core
}
