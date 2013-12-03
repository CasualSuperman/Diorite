package multiverse

import (
	"encoding/json"
	"math"
)

func (pt usuallyNumeric) MarshalJSON() ([]byte, error) {
	if math.IsNaN(float64(pt.Val)) {
		return json.Marshal(pt.Original)
	}
	return json.Marshal(pt.Val)
}
