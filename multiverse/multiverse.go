package multiverse

import (
	"github.com/CasualSuperman/phonetics/metaphone"
)

// Initialize the phonetics map for a constructed Multiverse.
func (m Multiverse) initialize() {
	markStandardSets(m.Sets)
	markExtendedSets(m.Sets)
}

func markStandardSets(sets []Set) {
	i := 0

	for !sets[i].Type.isTournamentLegal() {
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
		if sets[i].Type.isTournamentLegal() {
			sets[i].standardLegal = true

			if sets[i].Type == SetTypes.Core {
				coresInStandard--
			}
		}

		i++
	}
}

func markExtendedSets(sets []Set) {
	i := 0

	for !sets[i].Type.isTournamentLegal() {
		i++
	}

	mostRecentSet := sets[i]
	mostRecentRelease := mostRecentSet.Released
	releaseCutoff := mostRecentRelease.AddDate(-3, 0, 0)

	i = 0

	for sets[i].Released.After(releaseCutoff) {
		if sets[i].Type.isTournamentLegal() {
			sets[i].extendedLegal = true
		}

		i++
	}
}

// CardList is a list of cards.
type CardList []Card

// Add a card to the list, unless it's already in the list. Returns the card's position in the list.
func (c *CardList) Add(candidate Card) int {
	for i, card := range *c {
		if card.Name == candidate.Name {
			return i
		}
	}

	return c.AddNew(candidate)
}

// AddNew adds a new card to the list. Returns the card's position in the list.
func (c *CardList) AddNew(candidate Card) int {
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

// Len is the length of the CardList.
func (c CardList) Len() int {
	return len(c)
}

func (c *Card) scrub() {
	c.ascii = preventUnicode(c.Name)

	words := split(c.ascii)
	c.metaphones = make([]string, len(words))

	for j, str := range split(c.ascii) {
		c.metaphones[j] = metaphone.Encode(str)
	}
}
