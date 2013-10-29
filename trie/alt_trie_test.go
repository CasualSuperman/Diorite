package trie

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestAltAlt(t *testing.T) {
	trie := Alt()

	_, ok := trie.Get("Doesn't exist")

	if ok {
		t.Fail()
	}
}

func TestAltAdd(t *testing.T) {
	trie := Alt()

	if trie.Add("thing", 2) != nil {
		t.Fail()
	}
	if trie.Add("other", 3) != nil {
		t.Fail()
	}
}

func TestAltDuplicateAdd(t *testing.T) {
	trie := Alt()

	trie.Add("thing", 2)
	if trie.Add("thing", 7) == nil {
		t.Fail()
	}
}

func TestAltSubstringAdd(t *testing.T) {
	trie := Alt()

	trie.Add("sandlot", 1)
	err := trie.Add("sand", 7)

	if err != nil {
		t.Error("Substring wasn't added.")
	}

	val, ok := trie.Get("sand")

	if !ok {
		t.Fatal("Could not find substring")
	}
	if val.(int) != 7 {
		t.Error("Could not get value for substring")
	}
}

func TestAltSearch(t *testing.T) {
	trie := Alt()
	trie.Add("sand", 1)
	trie.Add("sandpaper", 2)
	trie.Add("sanity", 3)
	if x := trie.Search("san"); len(x) != 3 {
		t.Fail()
	}
	if x := trie.Search("sand"); len(x) != 2 {
		t.Fail()
	}
	results := trie.Search("sand")

	hasSand := false
	hasSandPaper := false

	for _, result := range results {
		if result == 1 {
			hasSand = true
		} else if result == 2 {
			hasSandPaper = true
		}
	}

	if !hasSandPaper || !hasSand {
		t.Fail()
	}
}

func TestAltRemove(t *testing.T) {
	trie := Alt()
	trie.Add("sand", 1)
	trie.Add("sandlot", 1)

	if x := trie.Remove("sand"); x != nil {
		t.Error("Couldn't remove sand.")
	} else if len(trie.Search("sand")) != 1 {
		fmt.Println(trie)
		t.Logf("%v\n", trie.Search("sand"))
		t.Error("Sand wasn't removed.")
	}

	trie.Add("sanity", 1)

	if x := trie.Remove("sandlot"); x != nil {
		t.Error("Couldn't remove sandlot.")
	} else if len(trie.Search("s")) != 1 {
		t.Logf("%v\n", trie.Search("s"))
		t.Error("Sandlot wasn't removed.")
	}

	if len(trie.children[0].branch.children[0].branch.children[0].branch.children) != 1 {
		t.Logf("%v\n", trie.children[0].branch.children[0].branch.children[0].branch.children)
		t.Error("Subtree wasn't deleted properly.")
	}
}

func BenchmarkAltAdd(b *testing.B) {
	list, err := ioutil.ReadFile("/usr/share/dict/words")

	if err != nil {
		b.Error(err)
	}

	words := strings.Split(string(list), "\n")
	t := Alt()

	b.ResetTimer()
	for i, j := 0, 0; i < b.N; i, j = i+1, j+1 {
		if j >= len(words) {
			j = 0
		}
		t.Add(words[j], true)
	}

}

func BenchmarkAltGet(b *testing.B) {
	list, err := ioutil.ReadFile("/usr/share/dict/words")

	if err != nil {
		b.Error(err)
	}

	words := strings.Split(string(list), "\n")
	t := Alt()

	for _, word := range words {
		t.Add(word, true)
	}

	b.ResetTimer()
	for i, j := 0, 0; i < b.N; i, j = i+1, j+1 {
		if j >= len(words) {
			j = 0
		}
		t.Get(words[j])
	}

}

func BenchmarkAltSearch(b *testing.B) {
	list, err := ioutil.ReadFile("/usr/share/dict/words")

	if err != nil {
		b.Error(err)
	}

	words := strings.Split(string(list), "\n")
	t := Alt()

	for _, word := range words {
		t.Add(word, true)
	}

	b.ResetTimer()
	for i, j := 0, 0; i < b.N; i, j = i+1, j+1 {
		if j >= len(words) {
			j = 0
		}
		t.Search(words[j])
	}

}
