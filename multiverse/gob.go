package multiverse

import (
	"encoding/gob"
	"io"
	"time"

	"code.google.com/p/lzma"
)

type gobMutiverse struct {
	Sets     []*Set
	CardList []gobCard
	Modified time.Time
}

type gobCard struct {
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

	Prints []gobPrinting

	Restricted, Banned []string
}

type gobPrinting struct {
	ID     MultiverseID
	Set    int
	Rarity Rarity
}

func (g *gobCard) card(sets []*Set) Card {
	var c Card

	c.Name = g.Name
	c.Cmc = g.Cmc
	c.Cost = g.Cost
	c.Colors = g.Colors

	c.Supertypes = g.Supertypes
	c.Types = g.Types

	c.Text = g.Text
	c.Flavor = g.Flavor

	c.Artist = g.Artist
	c.Number = g.Number

	c.Power = g.Power
	c.Toughness = g.Toughness

	c.Rulings = g.Rulings

	c.Printings = make([]Printing, len(g.Prints))

	for j, printing := range g.Prints {
		c.Printings[j].Rarity = printing.Rarity
		c.Printings[j].ID = printing.ID
		c.Printings[j].Set = sets[printing.Set]
	}

	for _, format := range g.Restricted {
		for _, f := range Formats.List {
			if f.Name == format {
				c.Restricted = append(c.Restricted, f)
				break
			}
		}
	}

	for _, format := range g.Banned {
		for _, f := range Formats.List {
			if f.Name == format {
				c.Banned = append(c.Banned, f)
				break
			}
		}
	}

	return c
}

func (c *Card) gobCard(sets []*Set) gobCard {
	var g gobCard

	g.Name = c.Name
	g.Cmc = c.Cmc
	g.Cost = c.Cost
	g.Colors = c.Colors

	g.Supertypes = c.Supertypes
	g.Types = c.Types

	g.Text = c.Text
	g.Flavor = c.Flavor

	g.Artist = c.Artist
	g.Number = c.Number

	g.Power = c.Power
	g.Toughness = c.Toughness

	g.Rulings = c.Rulings

	g.Prints = make([]gobPrinting, len(c.Printings))

	for _, format := range c.Restricted {
		g.Restricted = append(g.Restricted, format.Name)
	}
	for _, format := range c.Banned {
		g.Banned = append(g.Banned, format.Name)
	}

	for i, set := range sets {
		for j, printing := range c.Printings {
			if printing.Set == set {
				g.Prints[j].Rarity = printing.Rarity
				g.Prints[j].ID = printing.ID
				g.Prints[j].Set = i
				break
			}
		}
	}

	return g
}

// Write the multiverse to the provided writer.
func (m Multiverse) Write(w io.Writer) error {
	return m.WriteCompressLevel(w, lzma.DefaultCompression)
}

// WriteCompressLevel writes the multiverse to the provided writer with the given level of compression.
func (m Multiverse) WriteCompressLevel(w io.Writer, compressionLevel int) error {
	rawCards := make([]gobCard, len(m.Cards))

	for i, card := range m.Cards {
		rawCards[i] = card.gobCard(m.Sets)
	}

	mEnc := gobMutiverse{
		m.Sets,
		rawCards,
		m.Modified,
	}

	lw := lzma.NewWriterLevel(w, compressionLevel)
	defer lw.Close()

	enc := gob.NewEncoder(lw)

	return enc.Encode(mEnc)
}

// Read the multiverse from the provided reader.
func Read(r io.Reader) (m Multiverse, err error) {
	var mDec gobMutiverse

	lr := lzma.NewReader(r)
	dec := gob.NewDecoder(lr)

	err = dec.Decode(&mDec)

	lr.Close()

	if err != nil {
		return
	}

	cardList := make(CardList, len(mDec.CardList))
	for i := range mDec.CardList {
		cardList[i] = mDec.CardList[i].card(mDec.Sets)
		cardList[i].scrub()
	}

	m = Multiverse{
		mDec.Sets,
		cardList,
		mDec.Modified,
	}

	m.initialize()

	return m, nil
}
