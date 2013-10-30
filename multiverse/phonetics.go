package multiverse

import (
	"strings"
	"sync"
	"unicode"

	"github.com/CasualSuperman/Diorite/trie"
	"github.com/CasualSuperman/phonetics"
	"github.com/CasualSuperman/sift3"
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

var unicodeLock sync.RWMutex
var unicodeCache = make(map[string]string)

func preventUnicode(name string) string {
	unicodeLock.RLock()

	if cached, ok := unicodeCache[name]; ok {
		unicodeLock.RUnlock()
		return cached
	}
	unicodeLock.RUnlock()

	oldName := name
	name = strings.ToLower(name)

	if cached, ok := unicodeCache[name]; ok {
		unicodeLock.Lock()
		unicodeCache[oldName] = cached
		return cached
	}

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

	unicodeLock.Lock()
	unicodeCache[oldName] = clean
	unicodeCache[name] = clean
	unicodeLock.Unlock()

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
	aggregator := newFuzzySearchList(count)

	searchPhrase = preventUnicode(searchPhrase)
	searchGrams2 := newNGram(searchPhrase, 2)
	searchGrams3 := newNGram(searchPhrase, 3)

	for _, searchTerm := range Split(searchPhrase) {
		for _, result := range m.Pronunciations.Search(getMetaphone(searchTerm)) {
			for _, cardIndex := range result.([]int) {
				name := preventUnicode(m.Cards.List[cardIndex].Name)

				bestMatch := 0
				bestLen := 0
				for _, word := range Split(name) {
					match := phonetics.DifferenceSoundex(word, searchTerm)
					if match > bestMatch {
						bestMatch = match
						bestLen = len(word)
					}
				}

				similarity := searchGrams2.Similarity(name)
				similarity += searchGrams3.Similarity(name)
				similarity *= len(name) * bestMatch * bestLen

				if strings.Contains(name, searchPhrase) {
					similarity *= 50
				}

				similarity /= sift3.Sift(searchPhrase, name) + 1

				aggregator.Add(cardIndex, similarity)
			}
		}

		for cardIndex, card := range m.Cards.List {
			for _, word := range Split(preventUnicode(card.Name)) {
				if sift3.Sift(word, searchTerm) <= len(searchTerm)/3 {

					name := preventUnicode(card.Name)
					similarity := searchGrams2.Similarity(name)
					similarity += searchGrams3.Similarity(name)
					similarity *= len(name) * phonetics.DifferenceSoundex(word, searchTerm)
					similarity /= sift3.Sift(searchPhrase, name) + 1

					aggregator.Add(cardIndex, similarity)
				}
			}
		}
	}


	if len(aggregator.data) < count {
		count = len(aggregator.data)
	}

	results := make([]*Card, count)

	for i, card := range aggregator.data {
		results[i] = m.Cards.List[card.index]
	}

	return results
}
