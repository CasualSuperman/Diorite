package multiverse

import (
	"time"
)

// block, classic, commander ?

func init() {
	Formats.List = []*Format{
		Formats.Standard,
		Formats.Extended,
		Formats.Modern,
		Formats.Vintage,
		Formats.Legacy,
	}
}

type Format struct {
	SetOk func(*Set) bool
	Name  string
}

// The deck formats we know about. (Standard, Extended, Modern, etc.)
var Formats = struct {
	Standard, Extended, Modern, Vintage, Legacy, Un *Format

	List []*Format
}{
	&Format{standardSetLegal, "standard"},
	&Format{extendedSetLegal, "extended"},
	&Format{modernSetLegal, "modern"},
	&Format{vintageSetLegal, "vintage"},
	&Format{legacySetLegal, "legacy"},
	&Format{unSet, "un"},

	nil,
}

func unSet(s *Set) bool {
	return s.Type == SetTypes.Un
}

func vintageSetLegal(s *Set) bool {
	return !unSet(s)
}

func legacySetLegal(s *Set) bool {
	return !unSet(s)
}

var firstModernSet = time.Date(2003, time.July, 28, 0, 0, 0, 0, time.UTC)

func modernSetLegal(s *Set) bool {
	return !unSet(s) && !s.Released.Before(firstModernSet) &&
		(s.Type == SetTypes.Core || s.Type == SetTypes.Expansion)
}

func extendedSetLegal(s *Set) bool {
	return true
}

func standardSetLegal(s *Set) bool {
	return true
}

func (f *Format) Ok(c *Card) (bool, error) {
	for _, format := range c.Restricted {
		if format == f {
			return false, nil
		}
	}
	for _, format := range c.Banned {
		if format == f {
			return false, nil
		}
	}
	for _, printing := range c.Printings {
		if f.SetOk(printing.Set) {
			return true, nil
		}
	}
	return false, nil
}
