package multiverse

import (
	"runtime"
	"strings"
	"sync"
	"unicode"

	"github.com/CasualSuperman/Diorite/trie"
	"github.com/CasualSuperman/phonetics/metaphone"
	"github.com/CasualSuperman/phonetics/ngram"
	"github.com/CasualSuperman/sift3"
)

func generatePhoneticsMaps(cards []scrubbedCard) trie.Trie {
	metaphoneMap := trie.Alt()

	for i, c := range cards {
		for _, word := range Split(c.Ascii) {
			if len(word) < 4 {
				continue
			}
			mtp := metaphone.Encode(word)

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

	m := metaphone.Encode(s)
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
			if r == ' ' || unicode.IsLetter(r) || r == '_' || r == '-' {
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
				for i >= 1 && f.data[i].similarity > f.data[i-1].similarity {
					f.data[i-1], f.data[i] = f.data[i], f.data[i-1]
					i--
				}
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

	for i := myLen - 2; i >= 0; i-- {
		if f.data[i].similarity < similarity {
			f.data[i+1] = f.data[i]
			f.data[i].index = index
			f.data[i].similarity = similarity
		} else {
			return
		}
	}
	if f.data[0].similarity < similarity {
		f.data[0].index = index
		f.data[0].similarity = similarity
	}
}

// FuzzyNameSearch searches for a card with a similar name to the searchPhrase, and returns count or less of the most likely results.
func (m Multiverse) FuzzyNameSearch(searchPhrase string, count int) []*Card {
	var done sync.WaitGroup
	aggregator := newFuzzySearchList(count)

	groups := runtime.GOMAXPROCS(-1)

	totalCards := m.Cards.List.Len()
	groupInterval := totalCards / groups

	searchPhrase = preventUnicode(searchPhrase)
	searchGrams2 := ngram.New(searchPhrase, 2)
	searchGrams3 := ngram.New(searchPhrase, 3)

	for _, searchTerm := range Split(searchPhrase) {
		if len(searchTerm) == 0 {
			continue
		}
		searchMetaphone := metaphone.Encode(searchTerm)
		done.Add(groups)

		for i := 0; i < groups; i++ {
			start := i * groupInterval
			end := start + groupInterval
			if i == groups-1 {
				end = totalCards
			}

			go func(searchTerm, searchMetaphone string, start, end int) {
				defer done.Done()
				cards := m.Cards.List[start:end]
				for cardIndex := range cards {
					card := cards[cardIndex]
					name := card.Ascii
					metaphones := card.Metaphones

					if name == searchPhrase {
						//	println("EXACT MATCH")
						aggregator.Add(cardIndex+start, int(^uint(0)>>1))
						continue
					}

					bestMatch := int(^uint(0) >> 1)

					for _, metaphone := range metaphones {
						if len(metaphone) == 0 {
							continue
						}

						match := int(sift3.SiftASCII(metaphone, searchMetaphone))

						if match < bestMatch {
							bestMatch = match
						}
					}

					similarity := float32(searchGrams2.Similarity(name))
					similarity += float32(searchGrams3.Similarity(name))
					similarity -= float32(bestMatch * 2)
					dist := sift3.SiftASCII(searchPhrase, name)
					similarity -= float32(bestMatch) * dist * dist

					if similarity > 0 {
						aggregator.Add(cardIndex+start, int(similarity))
					}
				}
			}(searchTerm, searchMetaphone, start, end)
		}
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
