package multiverse

import "testing"

func TestSplit(t *testing.T) {
	s := split("The quick brown fix jumps over the lazy dog.")
	expected := []string{"The", "quick", "brown", "fix", "jumps", "over", "the", "lazy", "dog."}

	for i, str := range s {
		if expected[i] != str {
			t.Error("Split failed.")
		}
	}
}
