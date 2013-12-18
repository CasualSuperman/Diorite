package main

import (
	"hash/fnv"
	"sort"

	m "github.com/CasualSuperman/Diorite/multiverse"
)

type formatList struct {
	Format             *m.Format
	Banned, Restricted map[string]bool
}

func makeCardList(b map[string]bool) []string {
	s := make([]string, len(b))
	i := 0
	for key := range b {
		s[i] = key
		i++
	}

	sort.StringSlice(s).Sort()

	return s
}

func clearBanlists() {
	for i := range multiverse.Cards {
		multiverse.Cards[i].Banned = nil
		multiverse.Cards[i].Restricted = nil
	}
}

func generateFormatsHash(formats []formatList) uint64 {
	hash := fnv.New64()

	for _, f := range formats {
		banned := makeCardList(f.Banned)
		restricted := makeCardList(f.Restricted)

		for _, name := range banned {
			hash.Write([]byte(name))
		}

		for _, name := range restricted {
			hash.Write([]byte(name))
		}
	}

	return hash.Sum64()
}
