package multiverse

import (
	"strings"

	"github.com/arbovm/levenshtein"
)

type resultsList struct {
	cards      []*Card
	searchTerm string
}

func (r *resultsList) Len() int {
	return len(r.cards)
}

func (r *resultsList) Swap(i, j int) {
	r.cards[i], r.cards[j] = r.cards[j], r.cards[i]
}

func (r *resultsList) Less(i, j int) bool {
	iDist := len(r.searchTerm) * len(r.searchTerm)
	jDist := len(r.searchTerm) * len(r.searchTerm)

	for _, word := range strings.Split(preventUnicode(r.cards[i].Name), " ") {
		maxI := len(r.searchTerm)
		if maxI > len(word) {
			maxI = len(word)
		}
		dist := levenshtein.Distance(r.searchTerm, word[0:maxI])
		if dist < iDist {
			iDist = dist
		}
	}

	for _, word := range strings.Split(preventUnicode(r.cards[j].Name), " ") {
		maxJ := len(r.searchTerm)
		if maxJ > len(word) {
			maxJ = len(word)
		}
		dist := levenshtein.Distance(r.searchTerm, word[0:maxJ])
		if dist < jDist {
			jDist = dist
		}
	}

	return iDist < jDist
}
