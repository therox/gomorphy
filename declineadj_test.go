package gomorphy

import "testing"

type adjCase struct {
	word    string
	c       Case
	n       Number
	g       Gender
	animate bool
	want    string
}

// TestDeclineAdj covers "красный" in the matrix of 3 genders × 6 cases × 2 numbers.
// Masculine singular accusative and plural accusative are additionally checked
// for both animate=true and animate=false.
func TestDeclineAdj(t *testing.T) {
	cases := []adjCase{
		// === masc.sing — for accs, animacy is set by the parameter. ===
		{"красный", Nominative, Singular, Masculine, false, "красный"},
		{"красный", Genitive, Singular, Masculine, false, "красного"},
		{"красный", Dative, Singular, Masculine, false, "красному"},
		{"красный", Accusative, Singular, Masculine, false, "красный"},
		{"красный", Accusative, Singular, Masculine, true, "красного"},
		{"красный", Instrumental, Singular, Masculine, false, "красным"},
		{"красный", Prepositional, Singular, Masculine, false, "красном"},

		// === femn.sing — accs is independent of animacy (form "красную"). ===
		{"красный", Nominative, Singular, Feminine, false, "красная"},
		{"красный", Genitive, Singular, Feminine, false, "красной"},
		{"красный", Dative, Singular, Feminine, false, "красной"},
		{"красный", Accusative, Singular, Feminine, false, "красную"},
		{"красный", Accusative, Singular, Feminine, true, "красную"},
		{"красный", Instrumental, Singular, Feminine, false, "красной"},
		{"красный", Prepositional, Singular, Feminine, false, "красной"},

		// === neut.sing — accs == nomn regardless of animacy. ===
		{"красный", Nominative, Singular, Neuter, false, "красное"},
		{"красный", Genitive, Singular, Neuter, false, "красного"},
		{"красный", Dative, Singular, Neuter, false, "красному"},
		{"красный", Accusative, Singular, Neuter, false, "красное"},
		{"красный", Accusative, Singular, Neuter, true, "красное"},
		{"красный", Instrumental, Singular, Neuter, false, "красным"},
		{"красный", Prepositional, Singular, Neuter, false, "красном"},

		// === plur via masculine: forms are common, accs is determined by animate. ===
		{"красный", Nominative, Plural, Masculine, false, "красные"},
		{"красный", Genitive, Plural, Masculine, false, "красных"},
		{"красный", Dative, Plural, Masculine, false, "красным"},
		{"красный", Accusative, Plural, Masculine, false, "красные"},
		{"красный", Accusative, Plural, Masculine, true, "красных"},
		{"красный", Instrumental, Plural, Masculine, false, "красными"},
		{"красный", Prepositional, Plural, Masculine, false, "красных"},

		// === plur via feminine — should yield the same forms (gender ignored). ===
		{"красный", Nominative, Plural, Feminine, false, "красные"},
		{"красный", Genitive, Plural, Feminine, false, "красных"},
		{"красный", Dative, Plural, Feminine, false, "красным"},
		{"красный", Accusative, Plural, Feminine, false, "красные"},
		{"красный", Accusative, Plural, Feminine, true, "красных"},
		{"красный", Instrumental, Plural, Feminine, false, "красными"},
		{"красный", Prepositional, Plural, Feminine, false, "красных"},

		// === plur via neuter — same forms here as well. ===
		{"красный", Nominative, Plural, Neuter, false, "красные"},
		{"красный", Genitive, Plural, Neuter, false, "красных"},
		{"красный", Dative, Plural, Neuter, false, "красным"},
		{"красный", Accusative, Plural, Neuter, false, "красные"},
		{"красный", Accusative, Plural, Neuter, true, "красных"},
		{"красный", Instrumental, Plural, Neuter, false, "красными"},
		{"красный", Prepositional, Plural, Neuter, false, "красных"},
	}

	for _, tc := range cases {
		name := tc.word + "/" + caseTag(tc.c) + "/" + numberTag(tc.n) + "/" + genderTag(tc.g)
		if tc.c == Accusative {
			if tc.animate {
				name += "/anim"
			} else {
				name += "/inan"
			}
		}
		t.Run(name, func(t *testing.T) {
			got, err := DeclineAdj(tc.word, tc.c, tc.n, tc.g, tc.animate)
			if err != nil {
				t.Fatalf("DeclineAdj(%q, %v, %v, %v, %v) error: %v",
					tc.word, tc.c, tc.n, tc.g, tc.animate, err)
			}
			if got != tc.want {
				t.Fatalf("DeclineAdj(%q, %v, %v, %v, %v) = %q, expected %q",
					tc.word, tc.c, tc.n, tc.g, tc.animate, got, tc.want)
			}
		})
	}
}

// TestDeclineAdjWordNotFound — a word outside the dictionary returns an error.
func TestDeclineAdjWordNotFound(t *testing.T) {
	_, err := DeclineAdj("оранжевый", Genitive, Singular, Masculine, false)
	if err == nil {
		t.Fatal("expected an error for an unknown adjective")
	}
}
