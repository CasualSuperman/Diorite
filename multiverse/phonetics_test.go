package multiverse

import (
	"testing"
)

/*
func TestNameIsResult(t *testing.T) {
	m := openTestingMultiverse()
	var wait sync.WaitGroup
	groups := runtime.GOMAXPROCS(-1)
	wait.Add(groups)

	totalCards := m.Cards.List.Len()
	groupInterval := totalCards / groups

	for i := 0; i < groups; i++ {
		start := i * groupInterval
		end := start + groupInterval
		if i == groups-1 {
			end = totalCards
		}
		go func(list []scrubbedCard) {
			defer wait.Done()
			for _, card := range list {
				if results := m.FuzzyNameSearch(card.Card.Name, 1); results[0].Name != card.Card.Name {
					t.Errorf("Searching for '%s' should return that card. Got: %s\n", card.Card.Name, results[0].Name)
				}
			}
		}(m.Cards.List[start:end])
	}
	wait.Wait()
}
*/

func TestPhoneticCorrection(t *testing.T) {
	m := getTestingMultiverse(t)

	if results := m.FuzzyNameSearch("asheok", 1); results[0].Name != "Ashiok, Nightmare Weaver" {
		t.Errorf("Searching for 'asheok' should return Ashiok, Nightmare Weaver. Got: %s\n", results[0].Name)
	}
}

func BenchmarkSingleWordFuzzyNameSearch(b *testing.B) {
	m := getTestingMultiverse(nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.FuzzyNameSearch("asheok", 15)
	}
}

func BenchmarkMultiplePartialWordFuzzyNameSearch(b *testing.B) {
	m := getTestingMultiverse(nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.FuzzyNameSearch("ava me", 15)
	}
}

func BenchmarkSingleLetterFuzzyNameSearch(b *testing.B) {
	m := getTestingMultiverse(nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.FuzzyNameSearch("d", 15)
	}
}

func BenchmarkFullNameFuzzyNameSearch(b *testing.B) {
	m := getTestingMultiverse(nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.FuzzyNameSearch("Might of Old Krosa", 15)
	}
}
