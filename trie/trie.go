package trie

// A Trie is similar to a Map, but mapTrieost operations are O(log n), and it allows for finding similar elements.
type Trie interface {
	Add(string, interface{}) error
	Get(string) (interface{}, bool)
	Search(string) []interface{}
	Remove(string) error
	Update(string, interface{}) error
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

type traversalStack []stackNode

// Iter returns a fully iterative implementation of a Trie, which is faster and uses less stack space.
func New() *trie {
	return &trie{nil, nil, false}
}

func (t *trie) getOrAddChildBranch(r byte) *trie {
	for _, child := range t.children {
		if child.letter == r {
			return child.branch
		}
	}

	branch := branch{r, New()}
	t.children = append(t.children, branch)
	return branch.branch
}

func (t *trie) getChild(r byte) int {
	for i, child := range t.children {
		if child.letter == r {
			return i
		}
	}
	return -1
}

func (t *trie) getChildBranch(r byte) *trie {
	for _, child := range t.children {
		if child.letter == r {
			return child.branch
		}
	}
	return nil
}

func (t *trie) removeChildIndex(i int) {
	last := len(t.children)-1
	if i < last {
		t.children[i], t.children[last] = t.children[last], t.children[i]
	}
	t.children = t.children[:last]
}

// Add an element to the Trie, mapped to the given value.
func (t *trie) Add(key string, val interface{}) error {
	if key == "" {
		return emptyKey()
	}

	root := t

	for len(key) > 0 {
		root = root.getOrAddChildBranch(key[0])
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
		root = root.getChildBranch(key[0])
		key = key[1:]

		if root == nil {
			return nil, false
		}
	}

	if root.validLeaf {
		return root.value, true
	}

	return nil, false
}

// Search the Trie for all keys starting with the key.
// A full listing of the Trie is possible using t.Search("")
func (t *trie) Search(key string) []interface{} {
	var inlineStack [16]stackNode
	var results []interface{}

	root := t

	for len(key) > 0 {
		root = root.getChildBranch(key[0])
		key = key[1:]
		if root == nil {
			return results
		}
	}

	tip := 0
	// The first item on the stack is the last character of our key.
	inlineStack[0] = stackNode{-1, root}
	branch := traversalStack(inlineStack[0:1])

	// To help visualize, this is a depth-first search.
	for tip >= 0 {
		// Move on to the next sibling of the last leaf we processed.
		branch[tip].index++

		// Check to see if we're out of children.
		if branch[tip].index >= len(branch[tip].leaf.children) {
			// We are, so add ourselves if we're a valid leaf.
			if branch[tip].leaf.validLeaf {
				results = append(results, branch[tip].leaf.value)
			}
			// This branch is completely done; remove it from the stack.
			branch = branch[:tip]
			tip--
			continue
		}

		// Not out of children, push the next one onto the stack.
		branch = append(branch, stackNode{
			// We start at -1 because the first thing we do is increment this.
			-1,
			branch[tip].leaf.children[branch[tip].index].branch,
		})
		tip++
	}
	return results
}

// Remove the key from the Trie.
// The Trie will compact itself if possible.
func (t *trie) Remove(key string) error {
	var inlineStack [16]stackNode
	tip := -1
	branch := traversalStack(inlineStack[0:0])
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
	}

	branch = append(branch, stackNode{
		0,
		root,
	})
	tip++

	if !branch[tip].leaf.validLeaf {
		return notFound(key)
	}

	branch[tip].leaf.value = nil
	branch[tip].leaf.validLeaf = false

	// Chop off our dead branches.
	for tip > 0 {
		// If this branch isn't a valid leaf and has no children.
		if !branch[tip].leaf.validLeaf && len(branch[tip].leaf.children) == 0 {
			// Find our parent.
			trim := branch[tip-1]
			// Remove us from our parent. (quickly if possible)
			if len(trim.leaf.children) == 1 {
				trim.leaf.children = nil
			} else {
				trim.leaf.removeChildIndex(trim.index)
			}
			// Become our parent and repeat this process.
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
		root = root.getChildBranch(key[0])
		key = key[1:]
		
		if root == nil {
			return notFound(key)
		}
	}

	if !root.validLeaf {
		return notFound(key)
	}

	root.value = val

	return nil
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
