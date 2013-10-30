package multiverse

import (
	"strings"
	"sync"
	"unicode"

	"github.com/CasualSuperman/Diorite/trie"
	"github.com/CasualSuperman/phonetics"
	"github.com/CasualSuperman/sift3"
)

func generatePhoneticsMaps(cards []scrubbedCard) trie.Trie {
	metaphoneMap := trie.Alt()

	for i, c := range cards {
		for _, word := range Split(c.Ascii) {
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

var phoneticsLock sync.RWMutex
var phoneticsCache = make(map[string]string)

func getMetaphone(s string) string {
	phoneticsLock.RLock()
	if cached, ok := phoneticsCache[s]; ok {
		phoneticsLock.RUnlock()
		return cached
	}
	phoneticsLock.RUnlock()

	m := phonetics.EncodeMetaphone(s)
	phoneticsLock.Lock()
	phoneticsCache[s] = m
	phoneticsLock.Unlock()
	return m
}

func preventUnicode(name string) string {
	name = strings.ToLower(name)

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
			if r == ' ' || unicode.IsLetter(r) || r == '_' {
				clean += string(r)
			}
		}
	}

	return clean
}

type fuzzySearchList struct {
	sync.Mutex
	data []similarityItem
}

type similarityItem struct {
	index      int
	similarity int
}

func newFuzzySearchList(count int) fuzzySearchList {
	t := fuzzySearchList{}
	t.data = make([]similarityItem, 0, count)
	return t
}

func (f *fuzzySearchList) Add(index int, similarity int) {
	f.Lock()
	defer f.Unlock()
	for i, item := range f.data {
		if item.index == index {
			if f.data[i].similarity < similarity {
				f.data[i].similarity = similarity
			}
			return
		}
	}

	myLen := len(f.data)

	if myLen < cap(f.data) {
		f.data = f.data[:myLen+1]
		f.data[myLen] = similarityItem{index, similarity}
		myLen++
	}

	for i := myLen - 1; i >= 0; i-- {
		if f.data[i].similarity < similarity {
			if i < myLen-1 {
				f.data[i+1] = f.data[i]
			}
			f.data[i].index = index
			f.data[i].similarity = similarity
		} else {
			return
		}
	}
}

// FuzzyNameSearch searches for a card with a similar name to the searchPhrase, and returns count or less of the most likely results.
func (m Multiverse) FuzzyNameSearch(searchPhrase string, count int) []*Card {
	var done sync.WaitGroup
	aggregator := newFuzzySearchList(count)

	searchPhrase = preventUnicode(searchPhrase)
	searchGrams2 := newNGram(searchPhrase, 2)
	searchGrams3 := newNGram(searchPhrase, 3)

	for _, searchTerm := range Split(searchPhrase) {
		if len(searchTerm) == 0 {
			continue
		}
		done.Add(2)

		go func(searchTerm string) {
			defer done.Done()
			for _, result := range m.Pronunciations.Search(getMetaphone(searchTerm)) {
				for _, cardIndex := range result.([]int) {
					name := m.Cards.List[cardIndex].Ascii

					bestMatch := 0
					for _, word := range Split(name) {
						if len(word) == 0 {
							continue
						}
						match := phonetics.DifferenceSoundex(word, searchTerm)
						if match > bestMatch {
							bestMatch = match
						}
					}

					if bestMatch == 0 {
						continue
					}

					similarity := float32(searchGrams2.Similarity(name))
					similarity += float32(searchGrams3.Similarity(name))
					similarity /= float32(len(name) + len(searchPhrase))
					similarity *= float32(bestMatch)
					similarity /= float32(sift3.Sift(searchPhrase, name) + 1)

					if similarity != 0 {
						aggregator.Add(cardIndex, int(similarity))
					}
				}
			}
		}(searchTerm)

		go func(searchTerm string) {
			defer done.Done()
			for cardIndex := range m.Cards.List {
				name := m.Cards.List[cardIndex].Ascii

				if name == searchPhrase {
					aggregator.Add(cardIndex, int(^uint(0)>>1))
					continue
				}

				bestMatch := 0

				for _, word := range Split(name) {
					if len(word) == 0 {
						continue
					}
					if sift3.Sift(word, searchTerm) <= len(searchTerm)/3 {
						match := phonetics.DifferenceSoundex(word, searchTerm)
						if match > bestMatch {
							bestMatch = match
						}
					}
				}

				if bestMatch == 0 {
					continue
				}

				similarity := float32(searchGrams2.Similarity(name))
				similarity += float32(searchGrams3.Similarity(name))
				similarity /= float32(len(name) + len(searchPhrase))
				similarity *= float32(bestMatch)
				similarity /= float32(sift3.Sift(searchPhrase, name) + 1)

				if similarity != 0 {
					aggregator.Add(cardIndex, int(similarity))
				}
			}
		}(searchTerm)
	}

	done.Wait()

	if len(aggregator.data) < count {
		count = len(aggregator.data)
	}

	results := make([]*Card, count)

	for i, card := range aggregator.data {
		results[i] = m.Cards.List[card.index].Card
	}

	return results
}
