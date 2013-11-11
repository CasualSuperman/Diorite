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

func (m Multiverse) Card(id int) Card {
	return m.Cards[id]
}

type CardList []Card

func (c *CardList) Add(candidate Card) int {
	for i, card := range *c {
		if card.Name == candidate.Name {
			return i
		}
	}

	candidate.scrub()

	i := len(*c)

	if cap(*c) > i {
		*c = (*c)[:i+1]
	} else {
		newC := make(CardList, i+1, (i+1)*2)
		copy(newC, *c)
		*c = newC
	}
	(*c)[i] = candidate

	return i
}

func (c CardList) Len() int {
	return len(c)
}

func (c *Card) scrub() {
	c.ascii = preventUnicode(c.Name)

	words := Split(c.ascii)
	c.metaphones = make([]string, len(words))

	for j, str := range Split(c.ascii) {
		c.metaphones[j] = metaphone.Encode(str)
	}
}
