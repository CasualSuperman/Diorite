package main

import (
	"encoding/gob"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/CasualSuperman/Diorite/trie"
	"github.com/dotCypress/phonetics"
)

type Multiverse struct {
	Sets           map[string]*set
	Cards          map[multiverseID]*card
	Printings      map[*card][]multiverseID
	Formats        map[string][]*set
	Pronunciations trie.Trie
	Modified       time.Time
}

func LoadMultiverse(r io.Reader) Multiverse {
	d := gob.NewDecoder(r)
	var container struct {
		Data map[string]jsonSet
		Date time.Time
	}
	d.Decode(&container)
	return CreateMultiverse(container.Data, container.Date)
}

func CreateMultiverse(json map[string]jsonSet, modified time.Time) Multiverse {
	m := Multiverse{
		make(map[string]*set),
		make(map[multiverseID]*card),
		make(map[*card][]multiverseID),
		make(map[string][]*set),
		trie.New(),
		modified,
	}

	sets := make([]set, len(json))

	i := 0
	for _, set := range json {
		sets[i] = SetFromJson(set)
		addCardsToMap(set, &m.Cards, &m.Printings)
		i++
	}

	cards := make([]*card, len(m.Printings))

	i = 0
	for card, _ := range m.Printings {
		cards[i] = card
		i++
	}

	m.Pronunciations = generatePhoneticsMaps(cards)

	return m
}

func addCardsToMap(set jsonSet, cards *map[multiverseID]*card, printings *map[*card][]multiverseID) {
	for _, card := range set.Cards {
		c := CardFromJson(&card)
		s := (*printings)[c]
		s = append(s, multiverseID(card.MultiverseId))
		(*printings)[c] = s
		(*cards)[multiverseID(card.MultiverseId)] = c
	}
}

func (m *Multiverse) SearchByName(name string) []*card {
	aggregator := make(map[string]*card)

	for _, key := range m.Pronunciations.Search(phonetics.EncodeMetaphone(name)) {
		c, _ := m.Pronunciations.Get(key)
		crds := c.([]*card)
		for _, crd := range crds {
			for _, word := range strings.Split(crd.Name, " ") {
				max := len(name)
				if max > len(word) {
					max = len(word)
				}
				prefix := word[0:max]
				if phonetics.DifferenceSoundex(name, prefix) >= 85 && levenshtein(name, prefix) <= len(name)/2+1 {
					aggregator[crd.Name] = crd
					break
				}
			}
		}
	}

	for crd := range m.Printings {
		for _, word := range strings.Split(strings.ToLower(crd.Name), " ") {
			maxLen := len(name)
			if maxLen > len(word) {
				maxLen = len(word)
			}
			if levenshtein(name, word[0:maxLen]) < len(name)/4 {
				aggregator[crd.Name] = crd
			}
		}
	}

	finalResults := make([]*card, len(aggregator))
	i := 0
	for _, crd := range aggregator {
		finalResults[i] = crd
		i++
	}

	toSort := resultsList{finalResults, strings.ToLower(name)}
	sort.Sort(&toSort)

	return finalResults
}

type resultsList struct {
	cards      []*card
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

	for _, word := range strings.Split(strings.ToLower(r.cards[i].Name), " ") {
		maxI := len(r.searchTerm)
		if maxI > len(word) {
			maxI = len(word)
		}
		dist := levenshtein(r.searchTerm, word[0:maxI])
		if dist < iDist {
			iDist = dist
		}
	}

	for _, word := range strings.Split(strings.ToLower(r.cards[j].Name), " ") {
		maxJ := len(r.searchTerm)
		if maxJ > len(word) {
			maxJ = len(word)
		}
		dist := levenshtein(r.searchTerm, word[0:maxJ])
		if dist < jDist {
			jDist = dist
		}
	}

	return iDist < jDist
}
