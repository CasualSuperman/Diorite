package multiverse

import (
	"math"
	"strconv"
)

func (pt usuallyNumeric) MarshalJSON() ([]byte, error) {
	str := pt.Original

	if !math.IsNaN(float64(pt.Val)) {
		str = strconv.FormatFloat(float64(pt.Val), 'f', -1, 32)
	}

	return []byte(strconv.QuoteToASCII(str)), nil
}

// JsonDecode allows a Format to be restored using the encoding/json package.
func (f *Format) UnmarshalJSON(name []byte) error {
	name = name[1:]
	name = name[:len(name)-1]
	for _, format := range Formats.List {
		if format.Name == string(name) {
			*f = *format
			return nil
		}
	}
	return unrecognizedFormatErr(string(name))
}

// JsonEncode allows a Format to be saveed using the encoding/json package.
func (f *Format) MarshalJSON() ([]byte, error) {
	return []byte(strconv.QuoteToASCII(f.Name)), nil
}
