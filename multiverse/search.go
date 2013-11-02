package multiverse

import (
	"runtime"
)

type Filter func(*Card) bool

func (c CardList) Filter(f Filter) CardList {
	cores := runtime.GOMAXPROCS(-1)
	sectionLen := c.Len() / cores

	cardChan := make(chan *Card, 16)
	doneChan := make(chan bool)
	list := CardList(make([]scrubbedCard, 0))

	for i := 0; i < cores; i++ {
		start := sectionLen * i
		end := start + sectionLen
		if i == cores {
			end = c.Len()
		}

		go func(start, end int) {
			for _, card := range c[start:end] {
				if f(card.Card) {
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
			list.Add(c)
		}
	}

	return list
}
