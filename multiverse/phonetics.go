package multiverse

import (
	"strings"

	"github.com/CasualSuperman/Diorite/trie"
	"github.com/dotCypress/phonetics"
)

func generatePhoneticsMaps(cards []*Card) *trie.Trie {
	metaphoneMap := trie.New()

	for i, c := range cards {
		name := preventUnicode(c.Name)
		for _, word := range strings.Split(name, " ") {
			if len(word) < 4 {
				continue
			}
			mtp := phonetics.EncodeMetaphone(word)

			others, ok := metaphoneMap.Get(mtp)
			if ok {
				slice := others.([]int)
				slice = append(slice, i)
				metaphoneMap.Remove(mtp)
				metaphoneMap.Add(mtp, slice)
			} else {
				metaphoneMap.Add(mtp, []int{i})
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
