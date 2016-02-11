package multiverse

import (
	"encoding/gob"
	"io"
	"time"

	"github.com/lxq/lzma"
)

type gobMutiverse struct {
	Sets     []Set
	Cards    []gobCard
	Modified time.Time
}

type gobCard struct {
	Name   string
	Cmc    float32
	Cost   string
	Colors ManaColor

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

	Prints []gobPrinting

	Restricted, Banned []*Format
}

type gobPrinting struct {
	ID     MultiverseID
	Set    int
	Rarity Rarity
}

func (g *gobCard) card(sets []Set) Card {
	var c Card

	c.Name = g.Name
	c.Cmc = g.Cmc
	c.Cost = g.Cost
	c.Colors = g.Colors

	c.Supertypes = g.Supertypes
	c.Types = g.Types
	c.Subtypes = g.Subtypes

	c.Text = g.Text
	c.Flavor = g.Flavor

	c.Artist = g.Artist
	c.Number = g.Number

	c.Power = g.Power
	c.Toughness = g.Toughness
	c.Loyalty = g.Loyalty

	c.Rulings = g.Rulings

	c.Printings = make([]Printing, len(g.Prints))

	for j, printing := range g.Prints {
		c.Printings[j].Rarity = printing.Rarity
		c.Printings[j].ID = printing.ID
		c.Printings[j].Set = &sets[printing.Set]
	}

	c.Restricted = g.Restricted
	c.Banned = g.Banned

	return c
}

func (c *Card) gobCard(sets []Set) gobCard {
	var g gobCard

	g.Name = c.Name
	g.Cmc = c.Cmc
	g.Cost = c.Cost
	g.Colors = c.Colors

	g.Supertypes = c.Supertypes
	g.Types = c.Types
	g.Subtypes = c.Subtypes

	g.Text = c.Text
	g.Flavor = c.Flavor

	g.Artist = c.Artist
	g.Number = c.Number

	g.Power = c.Power
	g.Toughness = c.Toughness
	g.Loyalty = c.Loyalty

	g.Rulings = c.Rulings

	g.Prints = make([]gobPrinting, len(c.Printings))

	g.Restricted = c.Restricted
	g.Banned = c.Banned

	for i := range sets {
		for j, printing := range c.Printings {
			if printing.Set == &sets[i] {
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

	cardList := make([]Card, len(mDec.Cards))
	for i := range mDec.Cards {
		cardList[i] = mDec.Cards[i].card(mDec.Sets)
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
