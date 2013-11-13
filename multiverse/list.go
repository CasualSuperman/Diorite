package multiverse

import "sort"

type CardList []*Card
type sortableCardList struct {
	*CardList
	less Comparison
}

func (c CardList) Len() int {
	return len(c)
}

// Sort the results based on the given comparison.
func (c *CardList) Sort(f Comparison) {
	sortable := sortableCardList{
		c,
		f,
	}

	sort.Sort(sortable)
}

func (c *CardList) addEnsureNonDuplicate(candidate *Card) {
	for _, card := range *c {
		if card == candidate {
			return
		}
	}
	c.add(candidate)
}

func (c *CardList) add(card *Card) {
	l := len(*c)
	if l < cap(*c) {
		*c = (*c)[:l+1]
	} else {
		newList := make(CardList, l+1, (l+1)*2)
		copy(newList, *c)
		*c = newList
	}
	(*c)[l] = card
}

func (c *CardList) trim() {
	if len(*c) < cap(*c) {
		newList := make(CardList, len(*c))
		copy(newList, *c)
		*c = newList
	}
}

func (s sortableCardList) Swap(i, j int) {
	(*s.CardList)[i], (*s.CardList)[j] = (*s.CardList)[j], (*s.CardList)[i]
}

func (s sortableCardList) Less(i, j int) bool {
	return s.less((*s.CardList)[i], (*s.CardList)[j])
}
