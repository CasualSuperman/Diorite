package multiverse

import (
	"github.com/CasualSuperman/phonetics/metaphone"
)

func (m Multiverse) Loaded() bool {
	return !m.Modified.IsZero()
}

// Initialize the phonetics map for a constructed Multiverse.
func (m Multiverse) initialize() {
	markStandardSets(m.Sets)
	markExtendedSets(m.Sets)
	for i := range m.Cards {
		m.Cards[i].scrub()
	}
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

func (c *Card) scrub() {
	c.ascii = preventUnicode(c.Name)

	words := split(c.ascii)
	c.metaphones = make([]string, len(words))

	for j, str := range split(c.ascii) {
		c.metaphones[j] = metaphone.Encode(str)
	}
}
