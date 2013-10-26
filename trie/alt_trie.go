package trie

import "errors"

type altBranch struct {
	letter rune
	branch *altTrie
}

type altTrie struct {
	value     interface{}
	validLeaf bool
	children  []altBranch
}

// Alt returns an alternate implementation of a Trie, which is slightly faster for searching. Useful for Tries that are created and then infrequently changed.
func Alt() *altTrie {
	return &altTrie{nil, false, nil}
}

func (t altTrie) getChild(r rune) int {
	for i, child := range t.children {
		if child.letter == r {
			return i
		}
	}
	return -1
}

// Add an element to the Trie, mapped to the given value.
func (t *altTrie) Add(key string, val interface{}) error {
	runes := []rune(key)
	exists := t.add(runes, val)

	if exists {
		return errors.New("key already exists")
	}

	return nil
}

func (t *altTrie) add(r []rune, val interface{}) bool {
	if len(r) == 0 {
		return false
	}

	if i := t.getChild(r[0]); i != -1 {
		if len(r) > 1 {
			return t.children[i].branch.add(r[1:], val)
		}
		if t.children[i].branch.validLeaf {
			return true
		}
		child := t.children[i].branch
		child.validLeaf = true
		child.value = val
	} else {
		if len(r) > 1 {
			child := altBranch{
				r[0],
				&altTrie{},
			}
			t.children = append(t.children, child)

			return child.branch.add(r[1:], val)
		}
		t.children = append(t.children, altBranch{
			r[0],
			&altTrie{val, true, nil},
		})
	}
	return false
}

// Get a value from the Trie.
// Uses a comma ok format.
func (t *altTrie) Get(key string) (interface{}, bool) {
	if len(key) == 0 {
		return nil, false
	}
	return t.get([]rune(key))
}

func (t *altTrie) get(key []rune) (interface{}, bool) {
	if len(key) == 0 {
		return t.value, t.validLeaf
	}
	if i := t.getChild(key[0]); i != -1 {
		return t.children[i].branch.get(key[1:])
	}
	return nil, false
}

// Search the Trie for all keys starting with the key.
// A full listing of the Trie is possible using t.Search("")
func (t *altTrie) Search(key string) []string {
	results := t.search([]rune(key))
	for i, result := range results {
		results[i] = key + result
	}
	return results
}

func (t *altTrie) search(key []rune) []string {
	if len(key) == 0 {
		var options []string
		for _, child := range t.children {
			for _, option := range child.branch.search(key) {
				options = append(options, string(child.letter)+option)
			}
		}
		if t.validLeaf {
			options = append(options, "")
		}
		return options
	}
	if i := t.getChild(key[0]); i != -1 {
		return t.children[i].branch.search(key[1:])
	}
	return make([]string, 0)
}

// Remove the key from the Trie.
// The Trie will compact itself if possible.
func (t *altTrie) Remove(key string) error {
	runes := []rune(key)

	if !t.remove(runes) {
		errors.New("key not in trie")
	}

	return nil
}

func (t *altTrie) remove(key []rune) bool {
	if len(key) == 1 {
		if i := t.getChild(key[0]); i != -1 {
			child := t.children[i].branch
			if len(child.children) == 0 {
				t.children = append(t.children[:i], t.children[i+1:]...)
			} else {
				child.validLeaf = false
				child.value = nil
			}
			return true
		}
		return false
	}

	if i := t.getChild(key[0]); i != -1 {
		child := t.children[i].branch
		ret := child.remove(key[1:])

		if !child.validLeaf && len(child.children) == 0 {
			t.children = append(t.children[:i], t.children[i+1:]...)
		}
		return ret
	}
	return false
}
