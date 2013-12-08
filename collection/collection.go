package collection

import (
	m "github.com/CasualSuperman/Diorite/multiverse"
)

type Collection []Entry

type Entry struct {
	ID       m.MultiverseID
	Count    int
	Specific bool
	Foil     bool
}
