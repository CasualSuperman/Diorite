package multiverse

import (
	"encoding/gob"
	"fmt"
	"io"
	"time"

	"code.google.com/p/lzma"
	"github.com/glenn-brown/skiplist"
)

type skipListElem struct {
	ID        MultiverseID
	CardIndex int32
}

type gobMutiverse struct {
	Sets     []*Set
	Cards    []skipListElem
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
}

type gobPrinting struct {
	ID  MultiverseID
	Set int
}

func (g *gobCard) card(sets []*Set) *Card {
	var c Card

	c.Name = g.Name
	c.Cmc = g.Cmc
	c.Cost = g.Cost
	c.Colors = g.Colors

	c.Supertypes = g.Supertypes
	c.Types = g.Types

	c.Rarity = g.Rarity

	c.Text = g.Text
	c.Flavor = g.Flavor

	c.Artist = g.Artist
	c.Number = g.Number

	c.Power = g.Power
	c.Toughness = g.Toughness

	c.Rulings = g.Rulings

	c.Printings = make([]Printing, len(g.Prints))

	for j, printing := range g.Prints {
		c.Printings[j].Set = sets[printing.Set]
	}

	return &c
}

func (c *Card) gobCard(sets []*Set) gobCard {
	var g gobCard

	g.Name = c.Name
	g.Cmc = c.Cmc
	g.Cost = c.Cost
	g.Colors = c.Colors

	g.Supertypes = c.Supertypes
	g.Types = c.Types

	g.Rarity = c.Rarity

	g.Text = c.Text
	g.Flavor = c.Flavor

	g.Artist = c.Artist
	g.Number = c.Number

	g.Power = c.Power
	g.Toughness = c.Toughness

	g.Rulings = c.Rulings

	g.Prints = make([]gobPrinting, len(c.Printings))

setLoop:
	for i, set := range sets {
		for j, printing := range c.Printings {
			if printing.Set == set {
				g.Prints[j].Set = i
				continue setLoop
			}
		}
	}

	return g
}

// Write the multiverse to the provided writer.
func (m Multiverse) Write(w io.Writer) error {
	encCards := make([]skipListElem, m.Cards.Printings.Len())
	rawCards := make([]gobCard, len(m.Cards.List))

	for i, node := 0, m.Cards.Printings.Front(); node != nil; i, node = i+1, node.Next() {
		encCards[i] = skipListElem{MultiverseID(node.Key().(int)), int32(node.Value.(int))}
	}

	for i, card := range m.Cards.List {
		rawCards[i] = card.Card.gobCard(m.Sets)
	}

	mEnc := gobMutiverse{
		m.Sets,
		encCards,
		rawCards,
		m.Modified,
	}

	lw := lzma.NewWriter(w)
	enc := gob.NewEncoder(lw)
	err := enc.Encode(mEnc)

	if err != nil {
		return err
	}

	lw.Close()

	return nil
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

	decCards := skiplist.New()

	for _, elem := range mDec.Cards {
		decCards.Insert(int(elem.ID), int(elem.CardIndex))
	}

	scrubbedCardList := make(CardList, len(mDec.CardList))
	for i := range mDec.CardList {
		scrubbedCardList[i].Card = mDec.CardList[i].card(mDec.Sets)
	}
	scrubbedCardList.scrub()

	decPronunciations := generatePhoneticsMaps(scrubbedCardList)

	var cards = struct {
		Printings *skiplist.T
		List      CardList
	}{decCards, scrubbedCardList}

	m = Multiverse{
		mDec.Sets,
		cards,
		decPronunciations,
		mDec.Modified,
	}

	return m, nil
}
