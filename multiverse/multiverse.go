package multiverse

import (
	"io"
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
	Formats        map[string][]*Set
	Pronunciations *trie.Trie
	Modified       time.Time
}

func Inflate(r io.Reader) Multiverse {
	om := inflateOnlineMultiverse(r)
	return Create(om.Sets, om.Modified)
}

func Create(json map[string]jsonSet, modified time.Time) Multiverse {
	m := Multiverse{
		make(map[string]*Set),
		skiplist.New(),
		make(map[string][]*Set),
		nil,
		modified,
	}

	for _, set := range json {
		m.Sets[set.Name] = SetFromJson(set)
		for _, card := range set.Cards {
			m.Cards.Insert(card.MultiverseId, CardFromJson(&card))
		}
	}

	var cards []*Card

uniqueCard:
	for card := m.Cards.Front(); card != nil; card = card.Next() {
		for _, c := range cards {
			if c == card.Value.(*Card) {
				continue uniqueCard
			}
		}
		cards = append(cards, card.Value.(*Card))
	}

	m.Pronunciations = generatePhoneticsMaps(cards)

	return m
}

func addCardsToMap(set jsonSet, cards *map[multiverseID]*Card, printings *map[*Card][]multiverseID) {
	for _, card := range set.Cards {
		c := CardFromJson(&card)
		s := (*printings)[c]
		s = append(s, multiverseID(card.MultiverseId))
		(*printings)[c] = s
		(*cards)[multiverseID(card.MultiverseId)] = c
	}
}

func (m *Multiverse) SearchByName(name string) []*Card {
	aggregator := make(map[string]*Card)

	for _, key := range m.Pronunciations.Search(phonetics.EncodeMetaphone(name)) {
		c, _ := m.Pronunciations.Get(key)
		crds := c.([]*Card)
		for _, crd := range crds {
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
