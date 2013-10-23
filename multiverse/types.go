package multiverse

import (
	"math"
	"strconv"
	"time"
)

type manaColor byte
type borderColor byte
type rarity byte
type multiverseID int32
type setType byte

// The colors of mana that exist in the Multiverse.
var ManaColor = struct {
	White, Blue, Black, Red, Green manaColor
}{1, 2, 4, 8, 16}

// The borders that cards have.
var BorderColor = struct {
	White, Black, Silver borderColor
}{1, 2, 3}

// Rarities of cards.
var Rarity = struct {
	Common, Uncommon, Rare, Mythic, Basic, Special rarity
}{1, 2, 3, 4, 5, 6}

// Set types.
var SetType = struct {
	Core, Expansion, Reprint, Box, Un, FromTheVault, PremiumDeck, DuelDeck, Starter, Commander, Planechase, Archenemy setType
}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

const setReleaseFormat = "2006-01-02"

type Set struct {
	Name     string
	Code     string
	Released time.Time
	Border   borderColor
	Type     setType
	Block    string
	Cards    []multiverseID
}

type Card struct {
	Name   string
	Cmc    float32
	Cost   string
	Colors manaColor

	Supertypes, Types []string

	Rarity rarity

	Text   string
	Flavor string

	Artist string
	Number string

	Power, Toughness struct {
		Val      float32
		Original string
	}

	Rulings []ruling
}

type ruling struct {
	Date time.Time
	Text string
}

func copyCardFields(jc *jsonCard, c *Card) {
	c.Name = jc.Name
	c.Cmc = jc.Cmc
	c.Cost = jc.ManaCost

	for _, color := range jc.Colors {
		switch color {
		case "White":
			c.Colors |= ManaColor.White
		case "Blue":
			c.Colors |= ManaColor.Blue
		case "Black":
			c.Colors |= ManaColor.Black
		case "Red":
			c.Colors |= ManaColor.Red
		case "Green":
			c.Colors |= ManaColor.Green
		}
	}

	c.Supertypes = append(append(c.Supertypes, jc.Supertypes...), jc.Types...)
	c.Types = append(c.Types, jc.Subtypes...)

	switch jc.Rarity {
	case "Common":
		c.Rarity = Rarity.Common
	case "Uncommon":
		c.Rarity = Rarity.Uncommon
	case "Rare":
		c.Rarity = Rarity.Rare
	case "Mythic Rare":
		c.Rarity = Rarity.Mythic
	case "Special":
		c.Rarity = Rarity.Special
	case "Basic Land":
		c.Rarity = Rarity.Basic
	}

	c.Text = jc.Text
	c.Flavor = jc.Flavor
	c.Artist = jc.Artist
	c.Number = jc.Number

	c.Rulings = make([]ruling, len(jc.Rulings))

	power, err := strconv.ParseFloat(jc.Power, 32)

	if err == nil {
		c.Power.Val = float32(power)
	} else {
		c.Power.Val = float32(math.NaN())
		if jc.Power != "" {
			c.Power.Original = jc.Power
		}
	}

	toughness, err := strconv.ParseFloat(jc.Toughness, 32)

	if err == nil {
		c.Toughness.Val = float32(toughness)
	} else {
		c.Toughness.Val = float32(math.NaN())
		if jc.Toughness != "" {
			c.Toughness.Original = jc.Toughness
		}
	}
}

func SetFromJson(js jsonSet) *Set {
	t, _ := time.Parse(setReleaseFormat, js.ReleaseDate)
	var bColor borderColor
	var sType setType

	switch js.Border {
	case "black":
		bColor = BorderColor.Black
	case "white":
		bColor = BorderColor.White
	case "silver":
		bColor = BorderColor.Silver
	}

	switch js.Type {
	case "core":
		sType = SetType.Core
	case "expansion":
		sType = SetType.Expansion
	case "reprint":
		sType = SetType.Reprint
	case "box":
		sType = SetType.Box
	case "un":
		sType = SetType.Un
	case "from the vault":
		sType = SetType.FromTheVault
	case "premium deck":
		sType = SetType.PremiumDeck
	case "duel deck":
		sType = SetType.DuelDeck
	case "starter":
		sType = SetType.Starter
	case "commander":
		sType = SetType.Commander
	case "planechase":
		sType = SetType.Planechase
	case "archenemy":
		sType = SetType.Archenemy

	}

	ids := make([]multiverseID, len(js.Cards))

	i := 0
	for _, card := range js.Cards {
		ids[i] = multiverseID(card.MultiverseId)
		i++
	}

	return &Set{
		js.Name,
		js.Code,
		t,
		bColor,
		sType,
		js.Block,
		ids,
	}
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

func (c *Card) IsCreature() bool {
	for _, supertype := range c.Supertypes {
		if supertype == "Creature" {
			return true
		}
	}
	return false
}
