package multiverse

import (
	"io"
	"sort"
	"strings"
	"time"

	"github.com/CasualSuperman/Diorite/trie"
	"github.com/arbovm/levenshtein"
	"github.com/dotCypress/phonetics"
)

type Multiverse struct {
	Sets           map[string]*Set
	Cards          map[multiverseID]*Card
	Printings      map[*Card][]multiverseID
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
		make(map[multiverseID]*Card),
		make(map[*Card][]multiverseID),
		make(map[string][]*Set),
		nil,
		modified,
	}

	sets := make([]Set, len(json))

	i := 0
	for _, set := range json {
		sets[i] = SetFromJson(set)
		addCardsToMap(set, &m.Cards, &m.Printings)
		i++
	}

	cards := make([]*Card, len(m.Printings))

	i = 0
	for card, _ := range m.Printings {
		cards[i] = card
		i++
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

	for crd := range m.Printings {
		for _, word := range strings.Split(strings.ToLower(crd.Name), " ") {
			maxLen := len(name)
			if maxLen > len(word) {
				maxLen = len(word)
			}
			if levenshtein.Distance(name, word[0:maxLen]) < len(name)/4 {
				aggregator[crd.Name] = crd
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

	for _, word := range strings.Split(strings.ToLower(r.cards[i].Name), " ") {
		maxI := len(r.searchTerm)
		if maxI > len(word) {
			maxI = len(word)
		}
		dist := levenshtein.Distance(r.searchTerm, word[0:maxI])
		if dist < iDist {
			iDist = dist
		}
	}

	for _, word := range strings.Split(strings.ToLower(r.cards[j].Name), " ") {
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
