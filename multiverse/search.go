package multiverse

import (
	"fmt"
	"runtime"
	"sort"
	"strings"
)

// Ok makes Rarity usable as a Filter.
func (r Rarity) Ok(c *Card) (bool, error) {
	for _, printing := range c.Printings {
		if printing.Rarity == r {
			return true, nil
		}
	}
	return false, nil
}

// Ok makes ManaColor usable as a Filter.
func (m ManaColor) Ok(c *Card) (bool, error) {
	if m == ManaColors.Colorless {
		return c.Colors == 0, nil
	}
	return c.Colors&m != 0, nil
}

// Filter is a way to search through cards.
type Filter interface {
	Ok(*Card) (bool, error)
}

// Func is a generic type that allows a client to pass in any function that makes a boolean decision based on a card.
type Func func(*Card) bool

// Ok makes a Func usable as a Filter.
func (f Func) Ok(c *Card) (bool, error) {
	return f(c), nil
}

type searchResults []*Card

// Sort the results based on the given comparison.
func (s searchResults) Sort(f Comparison) searchResults {
	sortable := sortResults{
		s,
		f,
	}

	sort.Sort(sortable)

	return sortable.cards
}

// Search for cards that match the given conditions.
func (m Multiverse) Search(f Filter) (searchResults, error) {
	c := m.Cards
	cores := runtime.GOMAXPROCS(-1)
	sectionLen := c.Len() / cores

	cardChan := make(chan *Card, 16*cores)
	doneChan := make(chan bool)
	errChan := make(chan error)
	list := make(searchResults, 0, 1)

	for i := 0; i < cores; i++ {
		start := sectionLen * i
		end := start + sectionLen

		if i == cores-1 {
			end = c.Len()
		}

		go func(start, end int) {
			for j := range c[start:end] {
				ok, err := f.Ok(&c[j+start])
				if err != nil {
					errChan <- err
					return
				}
				if ok {
					cardChan <- &c[j+start]
				}
			}
			doneChan <- true
		}(start, end)
	}

	finishChans := func() {
		for cores > 0 {
			select {
			case <-cardChan:
			case <-errChan:
				cores--
			case <-doneChan:
				cores--
			}
		}

		close(cardChan)
		close(errChan)
		close(doneChan)
	}

	for cores > 0 {
		select {
		case <-doneChan:
			cores--
		case c := <-cardChan:
			appendNonDuplicateToCardList(&list, c)
		case err := <-errChan:
			cores--
			go finishChans()
			return nil, err
		}
	}

	close(cardChan)
	close(errChan)
	close(doneChan)

	for c := range cardChan {
		appendNonDuplicateToCardList(&list, c)

	}

	shortList := make(searchResults, len(list))

	copy(shortList, list)

	return shortList, nil
}

func appendToCardList(list *searchResults, c *Card) {
	for _, card := range *list {
		if card == c {
			return
		}
	}
	appendNonDuplicateToCardList(list, c)
}

func appendNonDuplicateToCardList(list *searchResults, c *Card) {
	l := len(*list)
	if l < cap(*list) {
		*list = (*list)[:l+1]
	} else {
		newList := make([]*Card, l+1, l*2)
		copy(newList, *list)
		*list = newList
	}
	(*list)[l] = c
}

// Not allows us to search for exclusive conditions.
type Not struct {
	Filter
}

// Ok makes Not usable as a Filter.
func (n Not) Ok(c *Card) (bool, error) {
	ok, err := n.Filter.Ok(c)
	return !ok, err
}

// And allows us to search for multiple conditions that must be true.
type And []Filter

// Ok makes And usable as a Filter.
func (a And) Ok(c *Card) (bool, error) {
	for _, f := range a {
		ok, err := f.Ok(c)
		if !ok || err != nil {
			return ok, err
		}
	}
	return true, nil
}

// Or allows us to search for multiple conditions that at least one of must be true.
// Performs short-circuit evaluation.
type Or []Filter

// Ok makes Or usable as a Filter.
func (o Or) Ok(c *Card) (bool, error) {
	for _, f := range o {
		ok, err := f.Ok(c)
		if ok || err != nil {
			return ok, err
		}
	}
	return false, nil
}

// Cond provides a way to search for non-builtin properties without resorting to a custom type.
type Cond map[string]interface{}

// Ok makes Cond usable as a Filter.
func (c Cond) Ok(card *Card) (bool, error) {
	for key, val := range c {
		switch key {
		case "color", "colors":
			same, err := handleColorSearch(card.Colors, val)
			if err != nil || !same {
				return same, err
			}
		case "cost":
			if val, ok := val.(string); ok {
				if strings.ToLower(card.Cost) != strings.ToLower(val) {
					return false, nil
				}
			} else {
				println("Unable to convert cost")
			}
		default:
			return false, fmt.Errorf("unsupported search method: %s", key)
		}
	}
	return true, nil
}

func handleColorSearch(cardColor ManaColor, val interface{}) (bool, error) {
	switch val := val.(type) {
	case ManaColor:
		if cardColor&val != 0 {
			return true, nil
		}
		return false, nil
	case string:
		switch strings.ToLower(val) {
		case "red":
			return cardColor&ManaColors.Red != 0, nil
		case "green":
			return cardColor&ManaColors.Green != 0, nil
		case "blue":
			return cardColor&ManaColors.Blue != 0, nil
		case "black":
			return cardColor&ManaColors.Black != 0, nil
		case "white":
			return cardColor&ManaColors.White != 0, nil
		case "colorless":
			return cardColor == 0, nil
		}
	}

	return false, fmt.Errorf("unexpected color type %T", val)
}
