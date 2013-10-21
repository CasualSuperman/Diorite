package multiverse

import (
	"sort"
	"strings"
	"time"

	"github.com/CasualSuperman/Diorite/trie"
	"github.com/arbovm/levenshtein"
	"github.com/dotCypress/phonetics"
	"github.com/glenn-brown/skiplist"
)

type Multiverse struct {
	Sets           map[string]*Set
	Cards          *skiplist.T
	cardList       []*Card
	Pronunciations *trie.Trie
	Modified       time.Time
}

func getCardIndex(cardList []*Card, cardName string) int {
	for i, card := range cardList {
		if card.Name == cardName {
			return i
		}
	}
	return -1
}

func Create(json map[string]jsonSet, modified time.Time) Multiverse {
	m := Multiverse{
		make(map[string]*Set),
		skiplist.New(),
		make([]*Card, 0),
		nil,
		modified,
	}

	for _, set := range json {
		m.Sets[set.Name] = SetFromJson(set)
		for _, card := range set.Cards {
			index := getCardIndex(m.cardList, card.Name)
			if index == -1 {
				index = len(m.cardList)
				c := new(Card)
				copyCardFields(&card, c)
				m.cardList = append(m.cardList, c)
			}
			m.Cards.Insert(card.MultiverseId, index)
		}
	}

	m.Pronunciations = generatePhoneticsMaps(m.cardList)

	return m
}

func (m *Multiverse) SearchByName(name string) []*Card {
	aggregator := make(map[string]*Card)

	for _, key := range m.Pronunciations.Search(phonetics.EncodeMetaphone(name)) {
		c, _ := m.Pronunciations.Get(key)
		crds := c.([]int)
		for _, i := range crds {
			crd := m.cardList[i]
			for _, word := range strings.Split(crd.Name, " ") {
				max := len(name)
				if max > len(word) {
					max = len(word)
				}
				prefix := word[0:max]
				if phonetics.DifferenceSoundex(name, prefix) >= 85 && levenshtein.Distance(name, prefix) <= len(name)/2+1 {
					aggregator[crd.Name] = crd
					break
				}
			}
		}
	}

	for crdName, crd := range processedCards.cards {
		for _, word := range strings.Split(strings.ToLower(crdName), " ") {
			maxLen := len(name)
			if maxLen > len(word) {
				maxLen = len(word)
			}
			if levenshtein.Distance(name, word[0:maxLen]) < len(name)/4 {
				aggregator[crdName] = crd
			}
		}
	}

	finalResults := make([]*Card, len(aggregator))
	i := 0
	for _, crd := range aggregator {
		finalResults[i] = crd
		i++
	}

	toSort := resultsList{finalResults, strings.ToLower(name)}
	sort.Sort(&toSort)

	return finalResults
}
