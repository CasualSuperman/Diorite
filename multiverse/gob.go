package multiverse

import (
	"encoding/gob"
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
	Sets     map[string]*Set
	Cards    []skipListElem
	CardList []*Card
	Modified time.Time
}

// Write the multiverse to the provided writer.
func (m Multiverse) Write(w io.Writer) error {
	encCards := make([]skipListElem, m.Cards.Printings.Len())
	rawCards := make([]*Card, len(m.Cards.List))

	for i, node := 0, m.Cards.Printings.Front(); node != nil; i, node = i+1, node.Next() {
		encCards[i] = skipListElem{MultiverseID(node.Key().(int)), int32(node.Value.(int))}
	}

	for i, card := range m.Cards.List {
		rawCards[i] = card.Card
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

	scrubbedCardList := scrubCards(mDec.CardList)
	decPronunciations := generatePhoneticsMaps(scrubbedCardList)

	var cards = struct {
		Printings *skiplist.T
		List      []scrubbedCard
	}{decCards, scrubbedCardList}

	m = Multiverse{
		mDec.Sets,
		cards,
		decPronunciations,
		mDec.Modified,
	}

	return m, nil
}
