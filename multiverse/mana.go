package multiverse

import (
	"fmt"
	"strconv"
	"strings"
)

type ManaCost []ManaType

func (m ManaCost) String() string {
	r := ""
	for _, s := range m {
		r += s.String()
	}
	return r
}

func (m ManaCost) Cmc() float32 {
	r := float32(0)
	for _, s := range m {
		r += s.Cmc()
	}
	return r
}

type manaCheck struct {
	symbol string
	result ManaType
}

var manaOptions = [][]manaCheck{
	baseManas,
	splitManas,
	varManas,
	phyrexianManas,
	unManas,
	snowManas,
}

var baseManas = []manaCheck{
	{"W", mana{ManaColors.White}},
	{"U", mana{ManaColors.Blue}},
	{"B", mana{ManaColors.Black}},
	{"R", mana{ManaColors.Red}},
	{"G", mana{ManaColors.Green}},
}

var splitManas = []manaCheck{
	{"W/U", slashMana{mana{ManaColors.White}, mana{ManaColors.Blue}}},
	{"W/B", slashMana{mana{ManaColors.White}, mana{ManaColors.Black}}},
	{"U/B", slashMana{mana{ManaColors.Blue}, mana{ManaColors.Black}}},
	{"U/R", slashMana{mana{ManaColors.Blue}, mana{ManaColors.Red}}},
	{"B/R", slashMana{mana{ManaColors.Black}, mana{ManaColors.Red}}},
	{"B/G", slashMana{mana{ManaColors.Black}, mana{ManaColors.Green}}},
	{"R/G", slashMana{mana{ManaColors.Red}, mana{ManaColors.Green}}},
	{"R/W", slashMana{mana{ManaColors.Red}, mana{ManaColors.White}}},
	{"G/W", slashMana{mana{ManaColors.Green}, mana{ManaColors.White}}},
	{"G/U", slashMana{mana{ManaColors.Green}, mana{ManaColors.Blue}}},
	{"2/W", slashMana{colorlessMana{"2", 2}, mana{ManaColors.White}}},
	{"2/U", slashMana{colorlessMana{"2", 2}, mana{ManaColors.Blue}}},
	{"2/B", slashMana{colorlessMana{"2", 2}, mana{ManaColors.Black}}},
	{"2/R", slashMana{colorlessMana{"2", 2}, mana{ManaColors.Red}}},
	{"2/G", slashMana{colorlessMana{"2", 2}, mana{ManaColors.Green}}},
}

var varManas = []manaCheck{
	{"X", colorlessMana{"X", 0}},
}

var phyrexianManas = []manaCheck{
	{"W/P", phyrexianMana{ManaColors.White}},
	{"U/P", phyrexianMana{ManaColors.Blue}},
	{"B/P", phyrexianMana{ManaColors.Black}},
	{"R/P", phyrexianMana{ManaColors.Red}},
	{"G/P", phyrexianMana{ManaColors.Green}},
}

var unManas = []manaCheck{
	{"hw", halfMana{mana{ManaColors.White}}},
	{"Y", colorlessMana{"Y", 0}},
	{"Z", colorlessMana{"Z", 0}},
}

var snowManas = []manaCheck{
	{"S", snowMana(true)},
}

func ParseCost(s string) (ManaCost, error) {
	cost := &s

	var m ManaCost

	for len(*cost) != 0 {
		circle, err := parseManaSymbol(cost)
		if err != nil {
			return nil, err
		}
		m = append(m, circle)
	}

	return m, nil
}

func parseManaSymbol(s *string) (ManaType, error) {
	cost := *s
	i := 0
	if cost[i] != '{' {
		return nil, fmt.Errorf("missing open bracket")
	}

	for i < len(cost) && cost[i] != '}' {
		i++
	}

	if i == len(cost) {
		return nil, fmt.Errorf("missing close bracket")
	}

	symbol := cost[1:i]
	*s = (*s)[i+1:]

	num, err := strconv.Atoi(symbol)

	// Colorless number?
	if err == nil {
		return colorlessMana{symbol, num}, nil
	}

	for _, manaType := range manaOptions {
		for _, option := range manaType {
			if option.symbol == symbol {
				return option.result, nil
			}
		}
	}

	return nil, fmt.Errorf("unknown mana format: %s", symbol)
}

func (m ManaColor) String() string {
	switch m {
	case ManaColors.White:
		return "W"
	case ManaColors.Blue:
		return "U"
	case ManaColors.Black:
		return "B"
	case ManaColors.Red:
		return "R"
	case ManaColors.Green:
		return "G"
	case ManaColors.Colorless:
		return "X"
	default:
		return "ERROR"
	}
}

// ManaType represents a single circle within a manacost.
// The MSB represents Phyrexian mana,
type ManaType interface {
	Cmc() float32
	Color() ManaColor
	String() string
}

type snowMana bool

func (m snowMana) Color() ManaColor {
	return ManaColors.Colorless
}

func (m snowMana) Cmc() float32 {
	return 1
}

func (m snowMana) String() string {
	return "{S}"
}

type mana struct {
	color ManaColor
}

func (m mana) Color() ManaColor {
	return m.color
}

func (m mana) Cmc() float32 {
	return 1
}

func (m mana) String() string {
	return "{" + m.color.String() + "}"
}

type halfMana struct {
	mana
}

func (h halfMana) Cmc() float32 {
	return 0.5
}

func (h halfMana) String() string {
	return "{h" + strings.ToLower(h.color.String()) + "}"
}

type slashMana struct {
	a, b ManaType
}

func (s slashMana) Color() ManaColor {
	return s.a.Color() & s.b.Color()
}

func (s slashMana) Cmc() float32 {
	cmc := s.a.Cmc()
	cmc2 := s.b.Cmc()
	if cmc > cmc2 {
		return cmc
	}
	return cmc2
}

func (s slashMana) String() string {
	a := s.a.String()
	b := s.b.String()
	return a[:len(a)-1] + "/" + b[1:]
}

type phyrexianMana mana

func (p phyrexianMana) Color() ManaColor {
	return mana(p).color
}

func (p phyrexianMana) Cmc() float32 {
	return 1
}

func (p phyrexianMana) String() string {
	return fmt.Sprintf("{P/%s}", p.color)
}

type colorlessMana struct {
	text string
	cmc  int
}

func (c colorlessMana) Color() ManaColor {
	return ManaColors.Colorless
}

func (c colorlessMana) Cmc() float32 {
	return float32(c.cmc)
}

func (c colorlessMana) String() string {
	return "{" + c.text + "}"
}
