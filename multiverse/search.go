package multiverse

import (
	"fmt"
	"runtime"
)

func (r Rarity) Ok(c *Card) (bool, error) {
	for _, printing := range c.Printings {
		if printing.Rarity == r {
			return true, nil
		}
	}
	return false, nil
}

func (m ManaColor) Ok(c *Card) (bool, error) {
	if m == ManaColors.Colorless {
		return c.Colors == 0, nil
	}
	return c.Colors&m != 0, nil
}

type Filter interface {
	Ok(*Card) (bool, error)
}

func (m Multiverse) Search(f Filter) ([]*Card, error) {
	c := m.Cards.List
	cores := runtime.GOMAXPROCS(-1)
	sectionLen := c.Len() / cores

	cardChan := make(chan *Card, cores*4)
	doneChan := make(chan bool)
	errChan := make(chan error)
	list := make([]*Card, 0, 1)

	for i := 0; i < cores; i++ {
		start := sectionLen * i
		end := start + sectionLen

		if i == cores-1 {
			end = c.Len()
		}

		go func(start, end int) {
			for _, card := range c[start:end] {
				ok, err := f.Ok(card.Card)
				if err != nil {
					errChan <- err
					return
				}
				if ok {
					cardChan <- card.Card
				}
			}
			doneChan <- true
		}(start, end)
	}

	for cores > 0 {
		select {
		case <-doneChan:
			cores--
		case c := <-cardChan:
			for _, card := range list {
				if card == c {
					continue
				}
			}
			l := len(list)
			if l < cap(list) {
				list = list[:l+1]
			} else {
				newList := make([]*Card, l+1, l*2)
				copy(newList, list)
				list = newList
			}
			list[l] = c
		case err := <-errChan:
			return nil, err
		}
	}

	shortList := make([]*Card, len(list))

	for i, card := range list {
		shortList[i] = card
	}

	return shortList, nil
}

type Not struct {
	f Filter
}

func (n Not) Ok(c *Card) (bool, error) {
	ok, err := n.f.Ok(c)
	return !ok, err
}

type And []Filter

func (a And) Ok(c *Card) (bool, error) {
	for _, f := range a {
		ok, err := f.Ok(c)
		if !ok || err != nil {
			return ok, err
		}
	}
	return true, nil
}

type Or []Filter

func (o Or) Ok(c *Card) (bool, error) {
	for _, f := range o {
		ok, err := f.Ok(c)
		if ok || err != nil {
			return ok, err
		}
	}
	return false, nil
}

type Cond map[string]interface{}

func (c Cond) Ok(card *Card) (bool, error) {
	for key, val := range c {
		switch key {
		case "color", "colors":
			same, err := handleColorSearch(card.Colors, val)
			if err != nil || !same {
				return same, err
			}
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
	}

	return false, fmt.Errorf("unexpected color type %T", val)
}
