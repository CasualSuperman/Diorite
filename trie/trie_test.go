package trie

import "testing"

func TestNew(t *testing.T) {
	trie := New()

	_, ok := trie.Get("Doesn't exist")

	if ok {
		t.Fail()
	}
}

func TestAdd(t *testing.T) {
	trie := New()

	if trie.Add("thing", 2) != nil {
		t.Fail()
	}
	if trie.Add("other", 3) != nil {
		t.Fail()
	}
}

func TestDuplicateAdd(t *testing.T) {
	trie := New()

	trie.Add("thing", 2)
	if trie.Add("thing", 7) == nil {
		t.Fail()
	}
}

func TestSubstringAdd(t *testing.T) {
	trie := New()

	trie.Add("sandlot", 1)
	err := trie.Add("sand", 7)

	if err != nil {
		t.Error("Substring wasn't added.")
	}

	val, ok := trie.Get("sand")

	if !ok {
		t.Error("Could not find substring")
	}
	if val.(int) != 7 {
		t.Error("Could not get value for substring")
	}
}

func TestSearch(t *testing.T) {
	trie := New()
	trie.Add("sand", 1)
	trie.Add("sandpaper", 1)
	trie.Add("sanity", 2)
	if x := trie.Search("san"); len(x) != 3 {
		t.Fail()
	}
	if x := trie.Search("sand"); len(x) != 2 {
		t.Fail()
	}
}

func TestRemove(t *testing.T) {
	trie := New()
	trie.Add("sand", 1)
	trie.Add("sandlot", 1)

	if x := trie.Remove("sand"); x != nil {
		t.Error("Couldn't remove sand.")
	} else if len(trie.Search("sand")) != 1 {
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

	if len(trie.children['s'].children['a'].children['n'].children) != 1 {
		t.Logf("%v\n", trie.children['s'].children['a'].children['n'].children)
		t.Error("Subtree wasn't deleted properly.")
	}
}
