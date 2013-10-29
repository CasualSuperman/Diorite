package multiverse

import (
	"strings"
	"unsafe"
)

type nGram struct {
	original string
	grams    []string
	size     int
}

func newNGram(phrase string, size int) nGram {
	var n = nGram{phrase, nil, size}
	if size == 2 || size == 3 {
		return n
	}
	words := strings.Split(phrase, " ")

	for _, word := range words {
		for i := size; i <= len(word); i++ {
			n.grams = append(n.grams, word[i-size:i])
		}
	}
	return n
}

func (n nGram) Similarity(phrase string) int {
	switch n.size {
	case 2:
		return n.twoSimilarity(phrase)
	case 3:
		return n.threeSimilarity(phrase)
	default:
		return n.nSimilarity(phrase)
	}
}

func (n nGram) twoSimilarity(phrase string) int {
	term := *((*[]byte)(unsafe.Pointer(&phrase)))
	original := *((*[]byte)(unsafe.Pointer(&n.original)))

	oLen := len(original)
	tLen := len(term)

	result := 0

	i := 0
	j := 0

	for i < oLen-1 && (original[i] == ' ' || original[i+1] == ' ') {
		i++
	}

	for j < tLen-1 && (term[j] == ' ' || term[j+1] == ' ') {
		j++
	}

	jStart := j

	for ; i < oLen-1; i++ {
		if original[i+1] == ' ' {
			i++
			continue
		}
		for j := jStart; j < tLen-1; j++ {
			if term[j+1] == ' ' {
				j++
				continue
			}

			if original[i] == term[j] && original[i+1] == term[j+1] {
				result++
			}
		}
	}

	return result * 2
}

func (n nGram) threeSimilarity(phrase string) int {
	term := *((*[]byte)(unsafe.Pointer(&phrase)))
	original := *((*[]byte)(unsafe.Pointer(&n.original)))

	oLen := len(original)
	tLen := len(term)

	result := 0

	i := 0
	j := 0

	for i < oLen-2 && (original[i] == ' ' || original[i+1] == ' ' || original[i+2] == ' ') {
		i++
	}

	for j < tLen-2 && (term[j] == ' ' || term[j+1] == ' ' || term[j+2] == ' ') {
		j++
	}

	jStart := j

	for ; i < oLen-2; i++ {
		if original[i+2] == ' ' {
			i += 2
			continue
		}
		for j := jStart; j < tLen-2; j++ {
			if term[j+2] == ' ' {
				j += 2
				continue
			}

			if original[i] == term[j] && original[i+1] == term[j+1] && original[i+2] == term[j+2] {
				result++
			}
		}
	}

	return result * 3
}

func (n nGram) nSimilarity(phrase string) int {
	result := 0

	i := 0

spaceLoop:
	for i < len(phrase)-n.size {
		for j := 0; j < n.size; j++ {
			if phrase[i+j] == ' ' {
				i++
				continue spaceLoop
			}
		}
		break
	}

	for ; i < len(phrase)-n.size; i++ {
		if phrase[i+n.size-1] == ' ' {
			i += n.size - 1
			continue
		}
		for _, myGram := range n.grams {
			if phrase[i:i+n.size] == myGram {
				result++
			}
		}
	}

	return result * n.size
}
