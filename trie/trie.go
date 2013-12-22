package trie

// A Trie is similar to a Map, but mapTrieost operations are O(log n), and it allows for finding similar elements.
type Trie interface {
	Add(string, interface{}) error
	Get(string) (interface{}, bool)
	Search(string) []interface{}
	Remove(string) error
	Update(string, interface{}) error
}

type (
	emptyKeyError bool
	keynotFoundError string
	duplicateKeyError string
)

func emptyKey() error {
	return emptyKeyError(true)
}

func notFound(key string) error {
	return keynotFoundError(key)
}

func duplicateKey(key string) error {
	return duplicateKeyError(key)
}

func (err emptyKeyError) Error() string {
	return "key empty"
}

func (err keynotFoundError) Error() string {
	return "key not in trie: '" + string(err) + "'"
}

func (err duplicateKeyError) Error() string {
	return "key already in trie: '" + string(err) + "'"
}

type trie struct {
	children  []branch
	value     interface{}
	validLeaf bool
}

type branch struct {
	letter byte
	branch *trie
}

type stackNode struct {
	index int
	leaf  *trie
}

// Iter returns a fully iterative implementation of a Trie, which is faster and uses less stack space.
func New() Trie {
	return &trie{nil, nil, false}
}

func (t *trie) getChild(r byte) int {
	for i, child := range t.children {
		if child.letter == r {
			return i
		}
	}
	return -1
}

// Add an element to the Trie, mapped to the given value.
func (t *trie) Add(key string, val interface{}) error {
	if key == "" {
		return emptyKey()
	}

	root := t

	for len(key) > 0 {
		i := root.getChild(key[0])

		if i == -1 {
			branch := branch{
				key[0],
				&trie{
					nil,
					nil,
					false,
				},
			}
			root.children = append(root.children, branch)
			root = branch.branch
		} else {
			root = root.children[i].branch
		}

		key = key[1:]
	}

	if root.validLeaf {
		return duplicateKey(key)
	}

	root.validLeaf = true
	root.value = val

	return nil
}

// Get a value from the Trie.
// Uses a comma ok format.
func (t *trie) Get(key string) (interface{}, bool) {
	if key == "" {
		return nil, false
	}

	root := t

	for len(key) > 0 {
		i := root.getChild(key[0])

		if i == -1 {
			return nil, false
		}

		root = root.children[i].branch
		key = key[1:]
	}

	if root.validLeaf {
		return root.value, true
	}

	return nil, false
}

// Search the Trie for all keys starting with the key.
// A full listing of the Trie is possible using t.Search("")
func (t *trie) Search(key string) []interface{} {
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
func (t *trie) Remove(key string) error {
	tip := -1
	branch := make([]stackNode, 0, 32)
	root := t

	for len(key) > 0 {
		i := root.getChild(key[0])

		if i == -1 {
			return notFound(key)
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
		return notFound(key)
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

	return nil
}

// Update the value of an existing element in the trie.
func (t *trie) Update(key string, val interface{}) error {
	if key == "" {
		return emptyKey()
	}

	root := t

	for len(key) > 0 {
		i := root.getChild(key[0])
		
		if i == -1 {
			return notFound(key)
		}

		root = root.children[i].branch
		key = key[1:]
	}

	if !root.validLeaf {
		return notFound(key)
	}

	root.value = val

	return nil
}
