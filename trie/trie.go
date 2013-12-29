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

func (t *trie) removeChild(c *trie) {
	last := len(t.children)-1

	for i, child := range t.children {
		if child.branch == c {
			if i < last {
				t.children[i], t.children[last] = t.children[last], t.children[i]
				break
			}
		}
	}
	t.children = t.children[:last]
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

	for i, l := 0, len(key); i < l; i++ {
		t = t.getOrAddChildBranch(key[i])
	}

	if t.validLeaf {
		return duplicateKey(key)
	}

	t.validLeaf = true
	t.value = val

	return nil
}

// Get a value from the Trie.
// Uses a comma ok format.
func (t *trie) Get(key string) (interface{}, bool) {
	if key == "" {
		return nil, false
	}

	for i, l := 0, len(key); i < l && t != nil; i++ {
		t = t.getChildBranch(key[i])
	}

	if t != nil && t.validLeaf {
		return t.value, true
	}

	return nil, false
}

// Search the Trie for all keys starting with the key.
// A full listing of the Trie is possible using t.Search("")
func (t *trie) Search(key string) []interface{} {
	var inlineStack [16]stackNode
	var results []interface{}

	for i, l := 0, len(key); i < l && t != nil; i++ {
		t = t.getChildBranch(key[i])
	}

	if t == nil {
		return results
	}

	tip := 0
	// The first item on the stack is the node positioned at the last character of our key.
	inlineStack[0] = stackNode{-1, t}
	stack := traversalStack(inlineStack[0:1])

	// To help visualize, this is a depth-first search.
	// We start at the node representing the end of the search key, and look at
	// its first child.  We put it on the stack, and proceed down the first
	// children of all these nodes until we hit a leaf.  We then check that
	// leaf to see if it is a validLeaf, and if it is, we put it onto the
	// results list.  After that, we increment the index into our children. If
	// this is greater than the number of children of the node, then we are
	// done with this node's children.  We then look at that node's validLeaf
	// status, and add it to the results if it's valid.  Then, we go on to that
	// node's nextSibling, etc.
	for tip >= 0 {
		// Move on to the next sibling of the last leaf we processed.
		stack[tip].index++

		// Check to see if we're out of children.
		if stack[tip].index >= len(stack[tip].leaf.children) {
			// We are, so add ourselves if we're a valid leaf.
			if stack[tip].leaf.validLeaf {
				results = append(results, stack[tip].leaf.value)
			}
			// This branch is completely done; remove it from the stack.
			stack = stack[:tip]
			tip--
			continue
		}

		next := stack[tip].leaf.children[stack[tip].index].branch

		// Avoid pushing leaves onto the stack.
		if len(next.children) > 0 {
			// Next node has children, push it onto the stack.
			stack = append(stack, stackNode{-1, next})
			tip++
		} else {
			// The next node doesn't have children, don't bother putting it on the stack.
			// Since we maintain a minimum tree, it will always be a validLeaf.
			results = append(results, next.value)
		}
	}
	return results
}

// Remove the key from the Trie.
// The Trie will compact itself if possible.
func (t *trie) Remove(key string) error {
	var inlineStack [16]*trie
	stack := inlineStack[0:0]

	// Identify the leaf associated with the key, and add every node we traverse to our stack.
	for i, l := 0, len(key); i < l && t != nil; i++ {
		stack = append(stack, t)
		t = t.getChildBranch(key[i])
	}

	if t == nil || !t.validLeaf {
		return notFound(key)
	}

	stack = append(stack, t)
	tip := len(key)

	stack[tip].value = nil
	stack[tip].validLeaf = false

	// Chop off our dead branches to free unneeded memory.
	for tip > 0 {
		// If this branch isn't a valid leaf and has no children.
		if stack[tip].validLeaf || len(stack[tip].children) > 0 {
			break
		}
		// Find our parent.
		trim := stack[tip-1]
		// Remove us from our parent. (quickly if possible)
		if len(trim.children) == 1 {
			trim.children = nil
		} else {
			trim.removeChild(stack[tip])
		}
		// Repeat this process with our parent.
		tip--
	}

	return nil
}

// Update the value of an existing element in the trie.
func (t *trie) Update(key string, val interface{}) error {
	if key == "" {
		return emptyKey()
	}

	for i, l := 0, len(key); i < l && t != nil; i++ {
		t = t.getChildBranch(key[i])
	}

	if t == nil || !t.validLeaf {
		return notFound(key)
	}

	t.value = val

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
