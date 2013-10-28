package multiverse

import "strings"

type nGram struct {
	grams []string
	size  int
}

func newNGram(phrase string, size int) nGram {
	var n = nGram{nil, size}
	words := strings.Split(phrase, " ")

	for _, word := range words {
		for i := size; i <= len(word); i++ {
			n.grams = append(n.grams, word[i-size:i])
		}
	}
	return n
}

func (n nGram) Similarity(phrase string) float32 {
	result := float32(0)
	m := newNGram(phrase, n.size)

	for _, myGram := range n.grams {
		for _, oGram := range m.grams {
			if myGram == oGram {
				result += float32(n.size)
			}
		}
	}

	return result
}
