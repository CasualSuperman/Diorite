package trie

import "errors"

type iterTrie struct {
	value     interface{}
	validLeaf bool
	children  []iterBranch
}

type iterBranch struct {
	letter byte
	branch *iterTrie
}

// Iter returns a fully iterative implementation of a Trie, which is faster and uses less stack space.
func Iter() *iterTrie {
	return &iterTrie{nil, false, nil}
}

func (t iterTrie) getChild(r byte) int {
	for i, child := range t.children {
		if child.letter == r {
			return i
		}
	}
	return -1
}

// Add an element to the Trie, mapped to the given value.
func (t *iterTrie) Add(key string, val interface{}) error {
	runes := []byte(key)
	if len(runes) == 0 {
		return errors.New("key empty")
	}
	exists := t.add(runes, val)

	if exists {
		return errors.New("key already exists")
	}

	return nil
}

func (t *iterTrie) add(r []byte, val interface{}) bool {
	root := t

	for {
		i := root.getChild(r[0])

		if len(r) > 1 {
			if i == -1 {
				branch := iterBranch{
					r[0],
					&iterTrie{
						nil,
						false,
						nil,
					},
				}
				root.children = append(root.children, branch)
				i = len(root.children) - 1
			}
			root = root.children[i].branch
			r = r[1:]
			continue
		}

		if i == -1 {
			leaf := iterBranch{
				r[0],
				&iterTrie{
					val,
					true,
					nil,
				},
			}
			root.children = append(root.children, leaf)
		} else {
			leaf := root.children[i].branch
			if leaf.validLeaf {
				break
			}

			leaf.validLeaf = true
			leaf.value = val
		}
		return false
	}

	return true
}

// Get a value from the Trie.
// Uses a comma ok format.
func (t *iterTrie) Get(key string) (interface{}, bool) {
	if len(key) == 0 {
		return nil, false
	}
	return t.get([]byte(key))
}

func (t *iterTrie) get(r []byte) (interface{}, bool) {
	root := t

	for {
		i := root.getChild(r[0])

		if len(r) > 1 {
			if i == -1 {
				break
			}
			root = root.children[i].branch
			r = r[1:]
			continue
		}

		if i == -1 {
			break
		}
		return root.children[i].branch.value, true
	}

	return nil, false
}

// Search the Trie for all keys starting with the key.
// A full listing of the Trie is possible using t.Search("")
func (t *iterTrie) Search(key string) []interface{} {
	return t.search([]byte(key))
}

type stackNode struct {
	index int
	leaf  *iterTrie
}

func (t *iterTrie) search(key []byte) []interface{} {
	root := t
	branch := make([]stackNode, 1, 32)
	results := make([]interface{}, 0)

	for len(key) > 0 {
		next := root.getChild(key[0])
		if next == -1 {
			return results
		}
		root = root.children[next].branch
		key = key[1:]
	}

	branch[0] = stackNode{-1, root}
	tip := 0

	for tip >= 0 {
		branch[tip].index++
		if branch[tip].index >= len(branch[tip].leaf.children) {
			if branch[tip].leaf.validLeaf {
				results = append(results, branch[tip].leaf.value)
			}
			branch = branch[:tip]
			tip--
			continue
		}

		next := branch[tip].leaf
		branch = append(branch, stackNode{
			-1,
			next.children[branch[tip].index].branch,
		})
		tip++
	}
	return results
}

// Remove the key from the Trie.
// The Trie will compact itself if possible.
func (t *iterTrie) Remove(key string) error {
	runes := []byte(key)

	if !t.remove(runes) {
		return errors.New("key not in trie")
	}

	return nil
}

func (t *iterTrie) remove(key []byte) bool {
	tip := -1
	branch := make([]stackNode, 0, 32)
	root := t

	for len(key) > 0 {
		i := root.getChild(key[0])
		if i == -1 {
			return false
		}
		branch = append(branch, stackNode{
			i,
			root,
		})
		root = root.children[i].branch
		key = key[1:]
		tip++

		if len(key) == 0 {
			branch = append(branch, stackNode{
				0,
				root,
			})
			tip++
		}
	}

	if branch[tip].leaf.validLeaf {
		branch[tip].leaf.value = nil
		branch[tip].leaf.validLeaf = false
	} else {
		return false
	}

	for tip > 0 {
		if !branch[tip].leaf.validLeaf && len(branch[tip].leaf.children) == 0 {
			trim := branch[tip-1]
			trim.leaf.children[trim.index] = trim.leaf.children[len(trim.leaf.children)-1]
			trim.leaf.children[len(trim.leaf.children)-1].branch = nil
			trim.leaf.children = trim.leaf.children[:len(trim.leaf.children)-1]
			if len(trim.leaf.children) == 0 {
				trim.leaf.children = nil
			}
			tip--
		} else {
			break
		}
	}

	return true
}

// Update the value of an existing element in the trie.
func (t *iterTrie) Update(key string, val interface{}) error {
	runes := []byte(key)
	ok := t.update(runes, val)
	if !ok {
		return errors.New("key is not in trie")
	}
	return nil
}

func (t *iterTrie) update(r []byte, val interface{}) bool {
	root := t

	for {
		i := root.getChild(r[0])

		if len(r) > 1 {
			if i == -1 {
				break
			}
			root = root.children[i].branch
			r = r[1:]
			continue
		}

		if i == -1 {
			break
		} else {
			leaf := root.children[i].branch
			if !leaf.validLeaf {
				break
			}

			leaf.value = val
		}
		return true
	}

	return false
}
