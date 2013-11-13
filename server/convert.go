package main

import (
	"math"
	"strconv"
	"time"

	m "github.com/CasualSuperman/Diorite/multiverse"
)

const setReleaseFormat = "2006-01-02"

func copyCardFields(jc *jsonCard, c *m.Card) {
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
}

func setFromJSON(js jsonSet) m.Set {
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
