package main

import (
	"runtime"
	"strings"
	"sync"

	"github.com/CasualSuperman/Diorite/trie"
	"github.com/dotCypress/phonetics"
)

func generatePhoneticsMaps(cards []*card) trie.Trie {
	metaphoneMap := trie.New()

	fixedCards := convertAllNames(cards)

	for _, c := range fixedCards {
		for _, word := range strings.Split(c.name, " ") {
			if len(word) < 4 {
				continue
			}
			mtp := phonetics.EncodeMetaphone(word)

			others, ok := metaphoneMap.Get(mtp)
			if ok {
				slice := others.([]*card)
				slice = append(slice, c.c)
				metaphoneMap.Remove(mtp)
				metaphoneMap.Add(mtp, slice)
			} else {
				metaphoneMap.Add(mtp, []*card{c.c})
			}
		}
	}

	return metaphoneMap
}

func preventUnicode(name string) string {
	for _, r := range name {
		if r > 128 {

		}
	}
	return strings.ToLower(name)
}

type namedCard struct {
	name string
	c    *card
}

func convertAllNames(cards []*card) []namedCard {
	var wg sync.WaitGroup
	workerCount := runtime.NumCPU()
	results := make([]namedCard, len(cards))

	pieceLen := len(cards) / workerCount

	for i := 0; i < workerCount; i++ {
		start := i * pieceLen
		end := start + pieceLen

		if i == workerCount-1 {
			end = len(cards)
		}

		wg.Add(1)
		go func(cards []*card, start int) {
			defer wg.Done()
			for i, c := range cards {
				results[i+start] = namedCard{
					preventUnicode(c.Name),
					c,
				}
			}
		}(cards[start:end], start)
	}

	wg.Wait()

	return results
}
