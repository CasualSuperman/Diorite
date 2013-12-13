package multiverse

import (
	"bytes"
	"runtime"
	"strings"
	"sync"
	"unicode"

	"github.com/CasualSuperman/phonetics/metaphone"
	"github.com/CasualSuperman/phonetics/ngram"
	"github.com/CasualSuperman/phonetics/sift3"
)

func preventUnicode(name string) string {
	var clean bytes.Buffer

	name = strings.ToLower(name)

	for _, r := range name {
		if r > 128 {
			switch r {
			case 'á', 'à', 'â':
				clean.WriteByte('a')
			case 'é':
				clean.WriteByte('e')
			case 'í':
				clean.WriteByte('i')
			case 'ö':
				clean.WriteByte('o')
			case 'û', 'ú':
				clean.WriteByte('u')

			case 'Æ', 'æ':
				clean.WriteString("ae")

			case '®':
				// We know this is an option but we're explicitly ignoring it.

			default:
			}
		} else {
			if r == ' ' || unicode.IsLetter(r) || r == '_' || r == '-' {
				clean.WriteRune(r)
			}
		}
	}

	return clean.String()
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
		if f.data[i].similarity < similarity || (f.data[i].similarity == similarity && f.data[i].index > index) {
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
func (m Multiverse) FuzzyNameSearch(searchPhrase string, count int) CardList {
	var done sync.WaitGroup
	aggregator := newFuzzySearchList(count)

	groups := runtime.GOMAXPROCS(-1)

	totalCards := len(m.Cards)
	groupInterval := totalCards / groups

	searchPhrase = preventUnicode(searchPhrase)
	searchGrams2 := ngram.New(searchPhrase, 2)
	searchGrams3 := ngram.New(searchPhrase, 3)

	for _, searchTerm := range split(searchPhrase) {
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
				cards := m.Cards[start:end]
				for cardIndex := range cards {
					card := cards[cardIndex]
					name := card.ascii
					metaphones := card.metaphones
					matchMod := float32(1.0)

					if name == searchPhrase {
						aggregator.Add(cardIndex+start, int(^uint(0)>>1))
						continue
					} else if strings.HasPrefix(name, searchPhrase) {
						matchMod = 1000
					}

					bestMatch := int(^uint(0) >> 1)

					for _, metaphone := range metaphones {
						if len(metaphone) == 0 {
							continue
						}

						match := int(sift3.SiftASCII(metaphone, searchMetaphone) / matchMod)

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

	results := make(CardList, count)

	for i, card := range aggregator.data {
		results[i] = &m.Cards[card.index]
	}

	return results
}
