package multiverse

// Sorts are built-in sorts available for use.
var Sorts = struct {
	Name, Cmc Comparison
}{
	nameSort,
	cmcSort,
}

// Comparison takes two cards and determines which is less than the other. See: sort.Interface.Less.
type Comparison func(*Card, *Card) bool

// Asc is a generated method that allows explicit specification that a sort will be ascending.
// It will follow the same sorting method as the base function.
func (c Comparison) Asc(a, b *Card) bool {
	return c(a, b)
}

// Desc is a generated method that allows a sort to be reversed without needing to write a wrapper.
func (c Comparison) Desc(a, b *Card) bool {
	return !c(a, b)
}

func nameSort(a, b *Card) bool {
	return a.Name < b.Name
}

func cmcSort(a, b *Card) bool {
	return a.Cmc < b.Cmc
}
