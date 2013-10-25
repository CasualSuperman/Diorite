package multiverse

import (
	"sort"
	"strings"
	"unicode"

	"github.com/CasualSuperman/Diorite/trie"
	"github.com/arbovm/levenshtein"
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

var seenRunes []rune

func preventUnicode(name string) string {
	clean := ""
	for _, r := range name {
		if r > 128 {
			switch r {
			case 'á', 'à', 'â':
				clean += "a"
			case 'é':
				clean += "e"
			case 'í':
				clean += "i"
			case 'ö':
				clean += "o"
			case 'û', 'ú':
				clean += "u"

			case 'Æ', 'æ':
				clean += "ae"

			case '®':
				// We know this is an option but we're explicitly ignoring it.

			default:
			}
		} else {
			if unicode.IsLetter(r) {
				clean += string(r)
			}
		}
	}
	return strings.ToLower(name)
}

type fuzzySearchList []struct {
	index      int
	similarity float32
}

func (f fuzzySearchList) Len() int {
	return len(f)
}

func (f fuzzySearchList) Less(i, j int) bool {
	return f[i].similarity > f[j].similarity
}

func (f fuzzySearchList) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (m Multiverse) FuzzyNameSearch(searchPhrase string, count int) []*Card {
	var aggregator fuzzySearchList
	searchPhrase = preventUnicode(searchPhrase)
	searchGrams2 := newNGram(searchPhrase, 2)
	searchGrams3 := newNGram(searchPhrase, 3)

	for _, searchTerm := range strings.Split(searchPhrase, " ") {
		for _, candidate := range m.Pronunciations.Search(phonetics.EncodeMetaphone(searchTerm)) {
			cardIndices, _ := m.Pronunciations.Get(candidate)
		cardLoop:
			for _, cardIndex := range cardIndices.([]int) {
				for _, i := range aggregator {
					if i.index == cardIndex {
						continue cardLoop
					}
				}
				name := preventUnicode(m.cardList[cardIndex].Name)
				similarity := searchGrams2.Similarity(name)
				similarity *= searchGrams3.Similarity(name)
				similarity /= float32(levenshtein.Distance(searchPhrase, name))
				similarity *= float32(len(name))

				if strings.Contains(name, searchPhrase) {
					similarity *= 10
				}

				var app = struct {
					index      int
					similarity float32
				}{
					cardIndex,
					similarity,
				}

				aggregator = append(aggregator, app)
			}
		}
	}

	sort.Sort(aggregator)
	if len(aggregator) < count {
		count = len(aggregator)
	}

	results := make([]*Card, count)

	for i := range results {
		results[i] = m.cardList[aggregator[i].index]
	}

	return results
}
