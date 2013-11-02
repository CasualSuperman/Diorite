package multiverse

import (
	"runtime"
)

type FilterFunc func(*Card) bool
type MapFunc func(*Card) interface{}

func (c CardList) Filter(f FilterFunc) CardList {
	cores := runtime.GOMAXPROCS(-1)
	sectionLen := c.Len() / cores

	cardChan := make(chan *Card, cores*4)
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

func (c CardList) Map(m MapFunc) []interface{} {
	cores := runtime.GOMAXPROCS(-1)
	sectionLen := c.Len() / cores

	cardChan := make(chan interface{}, cores*4)
	doneChan := make(chan bool)
	var list []interface{}

	for i := 0; i < cores; i++ {
		start := sectionLen * i
		end := start + sectionLen
		if i == cores {
			end = c.Len()
		}

		go func(start, end int) {
			for _, card := range c[start:end] {
				cardChan <- m(card.Card)
			}
			doneChan <- true
		}(start, end)
	}

	for cores > 0 {
		select {
		case <-doneChan:
			cores--
		case c := <-cardChan:
			l := len(list)
			if len(list) < cap(list) {
				list = list[:l]
				list[l] = c
			} else {
				nList := make([]interface{}, l+1, l*2)
				copy(nList, list)
				nList[l] = c
				list = nList
			}
		}
	}

	return list
}
