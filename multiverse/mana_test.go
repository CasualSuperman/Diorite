package multiverse

import "testing"

func TestParseCost(t *testing.T) {
	var tests = []struct {
		input string
		cmc   float32
	}{
		{"{1}", 1},
		{"{3}{G}{W}", 5},
		{"{X}{R}", 1},
		{"{2/G}{2/G}{2/G}", 6},
	}

	for _, test := range tests {
		r, err := ParseCost(test.input)
		t.Logf("Response string: %s, Original string: %s.\n", r.String(), test.input)
		t.Logf("Response cmc: %d, Original cmc: %d.\n", r.Cmc(), test.cmc)
		if err != nil {
			t.Log("Error:", err.Error())
			t.Fail()
		} else if test.input != r.String() {
			t.Fail()
		} else if test.cmc != r.Cmc() {
			t.Fail()
		}
	}
}
