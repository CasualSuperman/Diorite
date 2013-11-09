package multiverse

import (
	"github.com/CasualSuperman/phonetics/metaphone"
)

// Initialize the phonetics map for a constructed Multiverse.
func (m Multiverse) initialize() {
	markStandardSets(m.Sets)
	markExtendedSets(m.Sets)
}

func markStandardSets(sets []*Set) {
	i := 0

	for !sets[i].Type.IsTournamentLegal() {
		i++
	}

	mostRecentSet := sets[i]
	coresInStandard := 1

	mostRecentSetIsCoreSet := mostRecentSet.Type == SetTypes.Core

	if mostRecentSetIsCoreSet {
		coresInStandard++
	}

	i = 0

	for coresInStandard > 0 || sets[i].Type != SetTypes.Core {
		if sets[i].Type.IsTournamentLegal() {
			sets[i].standardLegal = true

			if sets[i].Type == SetTypes.Core {
				coresInStandard--
			}
		}

		i++
	}
}

func markExtendedSets(sets []*Set) {
	i := 0

	for !sets[i].Type.IsTournamentLegal() {
		i++
	}

	mostRecentSet := sets[i]
	mostRecentRelease := mostRecentSet.Released
	releaseCutoff := mostRecentRelease.AddDate(-3, 0, 0)

	i = 0

	for sets[i].Released.After(releaseCutoff) {
		if sets[i].Type.IsTournamentLegal() {
			sets[i].extendedLegal = true
		}

		i++
	}
}

func (m Multiverse) Card(id int) *Card {
	return m.Cards.List[id].Card
}

type CardList []scrubbedCard

func (c *CardList) Add(candidate *Card) int {
	for i, card := range *c {
		if card.Card == candidate {
			return i
		}
	}

	*c = append(*c, scrubbedCard{
		candidate,
		preventUnicode(candidate.Name),
		nil,
	})

	i := len(*c) - 1

	for _, str := range Split((*c)[i].Ascii) {
		(*c)[i].Metaphones = append((*c)[i].Metaphones, metaphone.Encode(str))
	}

	return i
}

func (c *CardList) Len() int {
	return len(*c)
}

type scrubbedCard struct {
	Card       *Card
	Ascii      string
	Metaphones []string
}

func (c CardList) scrub() {
	for i := range c {
		c[i].Ascii = preventUnicode(c[i].Card.Name)

		words := Split(c[i].Ascii)
		c[i].Metaphones = make([]string, len(words))

		for j, str := range Split(c[i].Ascii) {
			c[i].Metaphones[j] = metaphone.Encode(str)
		}
	}
}
