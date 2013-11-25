package multiverse

import (
	"bytes"
	"math"
	"strconv"
	"strings"
)

const cardWidth = 60

func (c Card) String() string {
	var b bytes.Buffer

	nameWidth := cardWidth - len(c.Cost) - 1
	b.WriteString(maxWidth(c.Name, nameWidth))
	for b.Len() < nameWidth {
		b.WriteByte(' ')
	}
	b.WriteString(c.Cost)
	b.WriteByte('\n')

	var typeLine bytes.Buffer

	typeLine.WriteString(c.Supertypes.String())

	if typeLine.Len() > 0 && c.Types != 0 {
		typeLine.WriteByte(' ')
	}

	typeLine.WriteString(c.Types.String())

	if len(c.Subtypes) > 0 {
		typeLine.WriteString(" - ")
		typeLine.WriteString(join(c.Subtypes))
	}

	b.Write(typeLine.Bytes())

	if len(c.Text) > 0 {
		b.WriteByte('\n')
		b.WriteString(maxWrap(c.Text, cardWidth))
	}

	if c.Is(Types.Creature) {
		b.WriteByte('\n')
		var pt bytes.Buffer

		if math.IsNaN(float64(c.Power.Val)) {
			pt.WriteString(c.Power.Original)
		} else {
			pt.WriteString(strconv.FormatFloat(float64(c.Power.Val), 'f', -1, 32))
		}
		pt.WriteByte('/')

		if math.IsNaN(float64(c.Toughness.Val)) {
			pt.WriteString(c.Toughness.Original)
		} else {
			pt.WriteString(strconv.FormatFloat(float64(c.Toughness.Val), 'f', -1, 32))
		}

		ptWidth := pt.Len()
		space := cardWidth - ptWidth - 1

		for space > 0 {
			b.WriteByte(' ')
			space--
		}
		b.Write(pt.Bytes())
	}

	return b.String()
}

func (s SuperType) String() string {
	var b bytes.Buffer

	if s&SuperTypes.Basic != 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString("Basic")
	}
	if s&SuperTypes.Elite != 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString("Elite")
	}
	if s&SuperTypes.Legendary != 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString("Legendary")
	}
	if s&SuperTypes.Ongoing != 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString("Ongoing")
	}
	if s&SuperTypes.Snow != 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString("Snow")
	}
	if s&SuperTypes.World != 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString("World")
	}

	return b.String()
}

func (t Type) String() string {
	var b bytes.Buffer
	if t&Types.Artifact != 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString("Artifact")
	}
	if t&Types.Creature != 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString("Creature")
	}
	if t&Types.Enchantment != 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString("Enchantment")
	}
	if t&Types.Instant != 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString("Instant")
	}
	if t&Types.Land != 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString("Land")
	}
	if t&Types.Planeswalker != 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString("Planeswalker")
	}
	if t&Types.Sorcery != 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString("Sorcery")
	}
	if t&Types.Tribal != 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString("Tribal")
	}

	return b.String()
}

func join(s []string) string {
	var b bytes.Buffer
	for i := range s {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(s[i])
	}
	return b.String()
}

func maxWidth(s string, width int) string {
	if len(s) <= width {
		return s
	}

	b := bytes.NewBuffer(nil)

	parts := split(s)
	i := 1
	for len(join(parts[:i])) < width {
		i++
	}

	b.WriteString(join(parts[:i-1]))
	l := b.Len()

	wordPartSize := width - l - 4

	if wordPartSize >= 3 {
		b.WriteByte(' ')
		b.WriteString(parts[i])
	}

	b.WriteString("...")

	return b.String()
}

func maxWrap(s string, width int) string {
	if len(s) <= width {
		return s
	}

	var b bytes.Buffer
	var line bytes.Buffer

	lines := strings.Split(s, "\n")

	for i, l := range lines {
		if i > 0 {
			b.WriteByte('\n')
		}

		parts := split(l)

		for i := range parts {
			if line.Len()+len(parts[i])+1 >= width {
				b.Write(line.Bytes())
				line.Reset()
				b.WriteByte('\n')
			} else if line.Len() > 0 {
				line.WriteByte(' ')
			}
			line.WriteString(parts[i])
		}

		b.Write(line.Bytes())
		line.Reset()
	}

	b.Write(line.Bytes())

	return b.String()
}
