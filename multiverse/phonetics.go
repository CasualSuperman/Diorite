package multiverse

import (
	"strings"

	"github.com/CasualSuperman/Diorite/trie"
	"github.com/dotCypress/phonetics"
)

func generatePhoneticsMaps(cards []*Card) *trie.Trie {
	metaphoneMap := trie.New()

	fixedCards := convertAllNames(cards)

	for _, c := range fixedCards {
		for _, word := range strings.Split(c.name, " ") {
			if len(word) < 4 {
				continue
			}
			mtp := phonetics.EncodeMetaphone(word)

			others, ok := metaphoneMap.Get(mtp)
			if ok {
				slice := others.([]*Card)
				slice = append(slice, c.c)
				metaphoneMap.Remove(mtp)
				metaphoneMap.Add(mtp, slice)
			} else {
				metaphoneMap.Add(mtp, []*Card{c.c})
			}
		}
	}

	return metaphoneMap
}

func preventUnicode(name string) string {
	for _, r := range name {
		if r > 128 {

		}
	}
	return strings.ToLower(name)
}

type namedCard struct {
	name string
	c    *Card
}

func convertAllNames(cards []*Card) []namedCard {
	results := make([]namedCard, len(cards))

	for i, card := range cards {
		results[i] = namedCard{
			preventUnicode(card.Name),
			card,
		}
	}

	return results
}
