package collection

import (
	"sort"

	m "github.com/CasualSuperman/Diorite/multiverse"
)

type Collection struct {
	entries []Entry
}

type Entry struct {
	ID   m.MultiverseID
	Sets []subEntry
}

type subEntry struct {
	Count    int
	Specific bool
	Foil     bool
}

func (c Collection) Len() int {
	return len(c.entries)
}

func (c Collection) Swap(i, j int) {
	c.entries[i], c.entries[j] = c.entries[j], c.entries[i]
}

func (c Collection) Less(i, j int) bool {
	return c.entries[i].ID < c.entries[j].ID
}

func (e *Entry) add(count int, specific, foil bool) {
	for i, set := range e.Sets {
		if set.Specific == specific && set.Foil == foil {
			e.Sets[i].Count += count
			return
		}
	}

	e.Sets = append(e.Sets, subEntry{count, specific, foil})
}

func (c *Collection) Add(id m.MultiverseID, count int, specific, foil bool) {
	i, found := c.search(id)
	if found {
		c.entries[i].add(count, specific, foil)
	} else {
		newEntries := make([]Entry, len(c.entries)+1)
		copy(newEntries, c.entries[0:i])
		newEntries[id] = Entry{
			id,
			[]subEntry{subEntry{count, specific, foil}},
		}
		copy(newEntries[i+1:], c.entries[i:])
		c.entries = newEntries
	}
}

func (c *Collection) Search(id m.MultiverseID) (i int, e Entry) {
	var found bool
	i, found = c.search(id)

	if !found {
		i = -1
	} else {
		e = c.entries[i]
	}

	return
}

func (c *Collection) search(id m.MultiverseID) (int, bool) {
	length := c.Len()

	i := sort.Search(length, func(index int) bool {
		return c.entries[index].ID <= id
	})

	if i == length {
		return i, false
	}

	return i, c.entries[i].ID == id
}
