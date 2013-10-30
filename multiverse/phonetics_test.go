package multiverse

import (
	"os"
	"testing"
)

func openTestingMultiverse() Multiverse {
	file, _ := os.Open("/home/rwertman/.diorite/multiverse.mtg")
	m, _ := Read(file)

	m.Initialize()
	return m
}

func TestPhoneticCorrection(t *testing.T) {
	m := openTestingMultiverse()

	if results := m.FuzzyNameSearch("asheok nightmare weaver", 1); results[0].Name != "Ashiok, Nightmare Weaver" {
		t.Errorf("Searching for 'asheok' should return Ashiok, Nightmare Weaver. Got: %s\n", results[0].Name)
	}
}

func BenchmarkSingleWordFuzzyNameSearch(b *testing.B) {
	m := openTestingMultiverse()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.FuzzyNameSearch("asheok", 15)
	}
}

func BenchmarkMultiplePartialWordFuzzyNameSearch(b *testing.B) {
	m := openTestingMultiverse()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.FuzzyNameSearch("ava me", 15)
	}
}

func BenchmarkSingleLetterFuzzyNameSearch(b *testing.B) {
	m := openTestingMultiverse()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.FuzzyNameSearch("d", 15)
	}
}

func BenchmarkFullNameFuzzyNameSearch(b *testing.B) {
	m := openTestingMultiverse()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.FuzzyNameSearch("Might of Old Krosa", 15)
	}
}
