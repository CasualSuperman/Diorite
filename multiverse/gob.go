package multiverse

import (
	"encoding/gob"
	"io"
	"time"

	"code.google.com/p/lzma"
	"github.com/glenn-brown/skiplist"
)

type skipListElem struct {
	ID        multiverseID
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
	encCards := make([]skipListElem, m.Cards.Len())

	for i, node := 0, m.Cards.Front(); node != nil; i, node = i+1, node.Next() {
		encCards[i] = skipListElem{multiverseID(node.Key().(int)), int32(node.Value.(int))}
	}

	mEnc := gobMutiverse{
		m.Sets,
		encCards,
		m.cardList,
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

	if err != nil {
		return
	}

	decCards := skiplist.New()

	for _, elem := range mDec.Cards {
		decCards.Insert(int32(elem.ID), int(elem.CardIndex))
	}

	decPronunciations := generatePhoneticsMaps(mDec.CardList)

	m = Multiverse{
		mDec.Sets,
		decCards,
		mDec.CardList,
		decPronunciations,
		mDec.Modified,
	}

	return m, nil
}
