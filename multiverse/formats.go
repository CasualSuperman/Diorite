package multiverse

import (
	"time"
)

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

// Format represents a sanctioned Magic the Gathering format.
type Format struct {
	SetOk func(*Set) bool
	Name  string
}

type unrecognizedFormatErr string

func (err unrecognizedFormatErr) Error() string {
	return "unrecognized format: " + string(err)
}

// GobDecode allows a Format to be restored using the encoding/gob package.
func (f *Format) GobDecode(name []byte) error {
	for _, format := range Formats.List {
		if format.Name == string(name) {
			*f = *format
			return nil
		}
	}
	return unrecognizedFormatErr(string(name))
}

// GobEncode allows a Format to be saveed using the encoding/gob package.
func (f *Format) GobEncode() ([]byte, error) {
	return []byte(f.Name), nil
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
	return s.extendedLegal
}

func standardSetLegal(s *Set) bool {
	return s.standardLegal
}

// Ok allows Formats to also act as a Filter for searching.
func (f Format) Ok(c *Card) (bool, error) {
	for _, format := range c.Restricted {
		if format.Name == f.Name {
			return false, nil
		}
	}
	for _, format := range c.Banned {
		if format.Name == f.Name {
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
