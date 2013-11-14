package main

import (
	"math"
       "sort"
	"strconv"
	"time"

	m "github.com/CasualSuperman/Diorite/multiverse"
)

const setReleaseFormat = "2006-01-02"

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
	c.Name = jc.Name
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

	c.Supertypes = append(append(c.Supertypes, jc.Supertypes...), jc.Types...)
	c.Types = append(c.Types, jc.Subtypes...)

	c.Text = jc.Text
	c.Flavor = jc.Flavor
	c.Artist = jc.Artist
	c.Number = jc.Number

	c.Rulings = make([]m.Ruling, len(jc.Rulings))

	for i, ruling := range jc.Rulings {
		c.Rulings[i].Text = ruling.Text
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
