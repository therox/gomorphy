package gomorphy

import "testing"

// TestGenderOf checks gender determination for all four categories and
// for pluralia tantum (where gender is undefined).
func TestGenderOf(t *testing.T) {
	cases := []struct {
		word string
		want Gender
	}{
		{"стол", Masculine},
		{"аппетит", Masculine},
		{"книга", Feminine},
		{"мать", Feminine},
		{"окно", Neuter},
		{"яблоко", Neuter},
		{"сирота", Common},
		// Pluralia tantum: the lemma has no gender tag — expect GenderUnknown.
		{"ножницы", GenderUnknown},
	}
	for _, tc := range cases {
		t.Run(tc.word, func(t *testing.T) {
			got, err := GenderOf(tc.word)
			if err != nil {
				t.Fatalf("GenderOf(%q) error: %v", tc.word, err)
			}
			if got != tc.want {
				t.Fatalf("GenderOf(%q) = %v, expected %v", tc.word, got, tc.want)
			}
		})
	}
}

// TestGenderOfWordNotFound — an unknown word → error.
func TestGenderOfWordNotFound(t *testing.T) {
	_, err := GenderOf("абракадабра")
	if err == nil {
		t.Fatal("expected an error for an unknown word")
	}
}
