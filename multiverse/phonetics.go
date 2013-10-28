package multiverse

import (
	"strings"
	"unicode"

	"github.com/CasualSuperman/Diorite/trie"
	"github.com/arbovm/levenshtein"
	"github.com/dotCypress/phonetics"
)

func generatePhoneticsMaps(cards []*Card) trie.Trie {
	metaphoneMap := trie.Alt()

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
			if r == ' ' || unicode.IsLetter(r) {
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

// FuzzyNameSearch searches for a card with a similar name to the searchPhrase, and returns count or less of the most likely results.
func (m Multiverse) FuzzyNameSearch(searchPhrase string, count int) []*Card {
	var aggregator = make(fuzzySearchList, 0, count)
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
				name := preventUnicode(m.Cards.List[cardIndex].Name)

				bestMatch := 0
				for _, word := range strings.Split(name, " ") {
					match := phonetics.DifferenceSoundex(word, searchTerm)
					if match > bestMatch {
						bestMatch = match
					}
				}

				similarity := searchGrams2.Similarity(name)
				similarity *= searchGrams3.Similarity(name)
				similarity *= float32(bestMatch) / 10.0
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

				if len(aggregator) < cap(aggregator) {
					i := len(aggregator) + 1
					aggregator = aggregator[:i]
					aggregator[i-1] = app
				} else {
					for i := count - 1; i >= 0; i-- {
						if aggregator[i].similarity < app.similarity {
							if i < count-1 {
								aggregator[i+1] = aggregator[i]
							}
							aggregator[i] = app
						} else {
							i = 0
						}
					}
				}
			}
		}
	levenshteinLoop:
		for cardIndex, card := range m.Cards.List {
			for _, word := range strings.Split(preventUnicode(card.Name), " ") {
				if levenshtein.Distance(word, searchTerm) <= len(searchTerm)/3 {

					name := preventUnicode(card.Name)
					similarity := searchGrams2.Similarity(name)
					similarity *= searchGrams3.Similarity(name)
					similarity *= float32(phonetics.DifferenceSoundex(word, searchTerm)) / 100.0
					similarity /= float32(levenshtein.Distance(searchPhrase, name))
					similarity *= float32(len(name))
					var app = struct {
						index      int
						similarity float32
					}{
						cardIndex,
						similarity,
					}

					for i, ci := range aggregator {
						if cardIndex == ci.index {
							if ci.similarity < similarity {
								aggregator[i] = app
							}
							continue levenshteinLoop
						}
					}

					if len(aggregator) < cap(aggregator) {
						i := len(aggregator) + 1
						aggregator = aggregator[:i]
						aggregator[i-1] = app
					} else {
						for i := count - 1; i >= 0; i-- {
							if aggregator[i].similarity < app.similarity {
								if i < count-1 {
									aggregator[i+1] = aggregator[i]
								}
								aggregator[i] = app
							} else {
								i = 0
							}
						}
					}
				}
			}
		}
	}

	if len(aggregator) < count {
		count = len(aggregator)
	}

	results := make([]*Card, count)

	for i, card := range aggregator {
		results[i] = m.Cards.List[card.index]
	}

	return results
}
