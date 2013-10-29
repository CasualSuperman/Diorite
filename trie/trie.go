package trie

import "errors"

// A Trie is similar to a Map, but mapTrieost operations are O(log n), and it allows for finding similar elements.
type Trie interface {
	Add(string, interface{}) error
	Get(string) (interface{}, bool)
	Search(string) []interface{}
	Remove(string) error
}

type mapTrie struct {
	value     interface{}
	validLeaf bool
	children  map[rune]*mapTrie
}

// New returns an initialized, empty Trie.
func New() *mapTrie {
	return &mapTrie{nil, false, make(map[rune]*mapTrie)}
}

// Add an element to the Trie, mapped to the given value.
func (t *mapTrie) Add(key string, val interface{}) error {
	runes := []rune(key)
	exists := t.add(runes, val)

	if exists {
		return errors.New("key already exists")
	}

	return nil
}

func (t mapTrie) add(r []rune, val interface{}) bool {
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
			child := mapTrie{
				nil,
				false,
				make(map[rune]*mapTrie),
			}
			t.children[r[0]] = &child

			return child.add(r[1:], val)
		}
		t.children[r[0]] = &mapTrie{
			val,
			true,
			make(map[rune]*mapTrie),
		}
	}
	return false
}

// Get a value from the Trie.
// Uses a comma ok format.
func (t *mapTrie) Get(key string) (interface{}, bool) {
	if len(key) == 0 {
		return nil, false
	}
	return t.get([]rune(key))
}

func (t mapTrie) get(key []rune) (interface{}, bool) {
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
func (t *mapTrie) Search(key string) []interface{} {
	results := t.search([]rune(key))
	if results == nil {
		results = make([]interface{}, 0)
	}
	return results
}

func (t mapTrie) search(key []rune) []interface{} {
	if len(key) == 0 {
		var options []interface{}
		for _, child := range t.children {
			options = append(options, child.search(key)...)
		}
		if t.validLeaf {
			options = append(options, t.value)
		}
		return options
	}
	if child, ok := t.children[key[0]]; ok {
		return child.search(key[1:])
	}
	return nil
}

// Remove the key from the Trie.
// The Trie will compact itself if possible.
func (t *mapTrie) Remove(key string) error {
	runes := []rune(key)

	if !t.remove(runes) {
		errors.New("key not in trie")
	}

	return nil
}

func (t mapTrie) remove(key []rune) bool {
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
