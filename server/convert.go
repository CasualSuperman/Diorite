package main

import (
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	m "github.com/CasualSuperman/Diorite/multiverse"
)

const setReleaseFormat = "2006-01-02"

func sanitize(s string) string {
	s = strings.Map(func(r rune) rune {
		if r == 0 {
			return -1
		}
		return r
	}, s)
	return s
}

// Convert to a Multiverse.
func (om onlineMultiverse) Convert() (mv m.Multiverse) {
	mv.Sets = make([]m.Set, len(om.Sets))
	mv.Modified = om.Modified

	i := 0

	for _, jSet := range om.Sets {
		mv.Sets[i] = jSet.Set()
		i++
	}

	s := setSorter{
		mv.Sets,
		releaseDateSort,
	}

	sort.Sort(s)

	cardIndexCache := make(map[string]int)

	for _, set := range om.Sets {
		var setIndex int

		for i, s := range mv.Sets {
			if s.Name == set.Name {
				setIndex = i
				break
			}
		}

		for _, jCard := range set.Cards {
			var rarity m.Rarity
			switch jCard.Rarity {
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
			default:
				panic("unknown rarity")
			}

			var printing = m.Printing{
				m.MultiverseID(jCard.MultiverseID),
				&mv.Sets[setIndex],
				rarity,
			}

			index, ok := cardIndexCache[jCard.Name]

			if ok {
				jCard := &mv.Cards[index]
				jCard.Printings = append(jCard.Printings, printing)
			} else {
				c := jCard.Card()
				c.Printings = []m.Printing{printing}
				i := len(mv.Cards)
				cardIndexCache[jCard.Name] = i

				if i < cap(mv.Cards) {
					mv.Cards = mv.Cards[:i+1]
				} else {
					newCards := make([]m.Card, i+1, (i+1)*2)
					copy(newCards, mv.Cards)
					mv.Cards = newCards
				}
				mv.Cards[i] = *c
			}
		}
	}

	newCards := make([]m.Card, len(mv.Cards))
	copy(newCards, mv.Cards)
	mv.Cards = newCards

	return
}

func (jc *jsonCard) Card() *m.Card {
	c := new(m.Card)
	c.Name = sanitize(jc.Name)
	c.Cmc = jc.Cmc
	c.Cost = jc.ManaCost

	for _, color := range jc.Colors {
		switch color {
		case "White":
			c.Colors |= m.ManaColors.White
		case "Blue":
			c.Colors |= m.ManaColors.Blue
		case "Black":
			c.Colors |= m.ManaColors.Black
		case "Red":
			c.Colors |= m.ManaColors.Red
		case "Green":
			c.Colors |= m.ManaColors.Green
		}
	}

	for _, SuperType := range jc.Supertypes {
		switch SuperType {
		case "Basic":
			c.Supertypes |= m.SuperTypes.Basic
		case "Elite":
			c.Supertypes |= m.SuperTypes.Elite
		case "Legendary":
			c.Supertypes |= m.SuperTypes.Legendary
		case "Ongoing":
			c.Supertypes |= m.SuperTypes.Ongoing
		case "Snow":
			c.Supertypes |= m.SuperTypes.Snow
		case "World":
			c.Supertypes |= m.SuperTypes.World
		}
	}

	for _, Type := range jc.Types {
		switch Type {
		case "Artifact":
			c.Types |= m.Types.Artifact
		case "Creature":
			c.Types |= m.Types.Creature
		case "Enchantment":
			c.Types |= m.Types.Enchantment
		case "Instant":
			c.Types |= m.Types.Instant
		case "Land":
			c.Types |= m.Types.Land
		case "Planeswalker":
			c.Types |= m.Types.Planeswalker
		case "Sorcery":
			c.Types |= m.Types.Sorcery
		case "Tribal":
			c.Types |= m.Types.Tribal
		}
	}

	c.Subtypes = jc.Subtypes

	c.Text = sanitize(jc.Text)
	c.Flavor = sanitize(jc.Flavor)
	c.Artist = sanitize(jc.Artist)
	c.Number = jc.Number

	c.Rulings = make([]m.Ruling, len(jc.Rulings))

	for i, ruling := range jc.Rulings {
		c.Rulings[i].Text = sanitize(ruling.Text)
		c.Rulings[i].Date, _ = time.Parse("2006-01-02", ruling.Date)
	}

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

	c.Loyalty = jc.Loyalty

	return c
}

func (js *jsonSet) Set() m.Set {
	s := m.Set{
		Name:  js.Name,
		Code:  js.Code,
		Block: js.Block,
	}
	s.Released, _ = time.Parse(setReleaseFormat, js.ReleaseDate)

	switch js.Border {
	case "black":
		s.Border = m.BorderColors.Black
	case "white":
		s.Border = m.BorderColors.White
	case "silver":
		s.Border = m.BorderColors.Silver
	}

	switch js.Type {
	case "core":
		s.Type = m.SetTypes.Core
	case "expansion":
		s.Type = m.SetTypes.Expansion
	case "reprint":
		s.Type = m.SetTypes.Reprint
	case "box":
		s.Type = m.SetTypes.Box
	case "un":
		s.Type = m.SetTypes.Un
	case "from the vault":
		s.Type = m.SetTypes.FromTheVault
	case "premium deck":
		s.Type = m.SetTypes.PremiumDeck
	case "duel deck":
		s.Type = m.SetTypes.DuelDeck
	case "starter":
		s.Type = m.SetTypes.Starter
	case "commander":
		s.Type = m.SetTypes.Commander
	case "planechase":
		s.Type = m.SetTypes.Planechase
	case "archenemy":
		s.Type = m.SetTypes.Archenemy
	}

	s.Cards = make([]m.MultiverseID, len(js.Cards))

	i := 0
	for _, card := range js.Cards {
		s.Cards[i] = m.MultiverseID(card.MultiverseID)
		i++
	}

	return s
}
