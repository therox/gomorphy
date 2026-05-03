package gomorphy

import "testing"

// TestPluralOf — nominative plural: ordinary words and suppletive pairs.
// For pluralia tantum the lemma is returned (it is already plural).
func TestPluralOf(t *testing.T) {
	cases := []struct {
		word string
		want string
	}{
		{"стол", "столы"},
		{"книга", "книги"},
		{"окно", "окна"},
		{"аппетит", "аппетиты"},
		{"яблоко", "яблоки"},
		// Suppletive pairs.
		{"человек", "люди"},
		{"ребёнок", "дети"},
		// Pluralia tantum — the lemma equals the plural.
		{"ножницы", "ножницы"},
	}
	for _, tc := range cases {
		t.Run(tc.word, func(t *testing.T) {
			got, err := PluralOf(tc.word)
			if err != nil {
				t.Fatalf("PluralOf(%q) error: %v", tc.word, err)
			}
			if got != tc.want {
				t.Fatalf("PluralOf(%q) = %q, expected %q", tc.word, got, tc.want)
			}
		})
	}
}

// TestPluralOfWordNotFound — an unknown word → error.
func TestPluralOfWordNotFound(t *testing.T) {
	_, err := PluralOf("абракадабра")
	if err == nil {
		t.Fatal("expected an error for an unknown word")
	}
}
