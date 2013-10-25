package trie

import "errors"

// A Trie is similar to a Map, but most operations are O(log n), and it allows for finding similar elements.
type Trie struct {
	value     interface{}
	validLeaf bool
	children  map[rune]*Trie
}

// New returns an initialized, empty Trie.
func New() *Trie {
	return &Trie{' ', false, make(map[rune]*Trie)}
}

// Add an element to the Trie, mapped to the given value.
func (t Trie) Add(key string, val interface{}) error {
	runes := []rune(key)
	exists := t.add(runes, val)

	if exists {
		return errors.New("key already exists")
	}

	return nil
}

func (t Trie) add(r []rune, val interface{}) bool {
	if len(r) == 0 {
		return false
	}

	if child, ok := t.children[r[0]]; ok {
		if len(r) > 1 {
			return child.add(r[1:], val)
		}
		if child.validLeaf {
			return true
		}
		child.validLeaf = true
		child.value = val
	} else {
		if len(r) > 1 {
			child := Trie{
				nil,
				false,
				make(map[rune]*Trie),
			}
			t.children[r[0]] = &child

			return child.add(r[1:], val)
		}
		t.children[r[0]] = &Trie{
			val,
			true,
			make(map[rune]*Trie),
		}
	}
	return false
}

// Get a value from the Trie.
// Uses a comma ok format.
func (t Trie) Get(key string) (interface{}, bool) {
	if len(key) == 0 {
		return nil, false
	}
	return t.get([]rune(key))
}

func (t Trie) get(key []rune) (interface{}, bool) {
	if len(key) == 0 {
		return t.value, t.validLeaf
	}
	if child, ok := t.children[key[0]]; ok {
		return child.get(key[1:])
	}
	return nil, false
}

// Search the Trie for all keys starting with the key.
// A full listing of the Trie is possible using t.Search("")
func (t Trie) Search(key string) []string {
	results := t.search([]rune(key))
	for i, result := range results {
		results[i] = key + result
	}
	return results
}

func (t Trie) search(key []rune) []string {
	if len(key) == 0 {
		var options []string
		for r, child := range t.children {
			for _, option := range child.search(key) {
				options = append(options, string(r)+option)
			}
		}
		if t.validLeaf {
			options = append(options, "")
		}
		return options
	}
	if child, ok := t.children[key[0]]; ok {
		return child.search(key[1:])
	}
	return make([]string, 0)
}

// Remove the key from the Trie.
// The Trie will compact itself if possible.
func (t Trie) Remove(key string) error {
	runes := []rune(key)

	if !t.remove(runes) {
		errors.New("key not in trie")
	}

	return nil
}

func (t Trie) remove(key []rune) bool {
	if len(key) == 1 {
		if child, ok := t.children[key[0]]; ok {
			if len(child.children) == 0 {
				delete(t.children, key[0])
			} else {
				child.validLeaf = false
				child.value = nil
			}
			return true
		}
		return false
	}

	if child, ok := t.children[key[0]]; ok {
		ret := child.remove(key[1:])

		if !child.validLeaf && len(child.children) == 0 {
			delete(t.children, key[0])
		}
		return ret
	}
	return false
}
