package gomorphy

import "testing"

// patrCase — one expectation row: input patronymic, case, expected form.
type patrCase struct {
	patr string
	c    Case
	want string
}

// TestPatronymicHeuristicMasc covers masculine patronymics in -ович, -евич,
// -ич, -ьич in six cases. The tests go straight through the heuristic,
// without the dictionary.
func TestPatronymicHeuristicMasc(t *testing.T) {
	cases := []patrCase{
		// -ович (Иванович) — stem "иванов".
		{"иванович", Nominative, "иванович"},
		{"иванович", Genitive, "ивановича"},
		{"иванович", Dative, "ивановичу"},
		{"иванович", Accusative, "ивановича"},
		{"иванович", Instrumental, "ивановичем"},
		{"иванович", Prepositional, "ивановиче"},

		// -евич (Сергеевич) — stem "сергеев".
		{"сергеевич", Nominative, "сергеевич"},
		{"сергеевич", Genitive, "сергеевича"},
		{"сергеевич", Dative, "сергеевичу"},
		{"сергеевич", Accusative, "сергеевича"},
		{"сергеевич", Instrumental, "сергеевичем"},
		{"сергеевич", Prepositional, "сергеевиче"},

		// -ич (Никитич) — stem "никит".
		{"никитич", Nominative, "никитич"},
		{"никитич", Genitive, "никитича"},
		{"никитич", Dative, "никитичу"},
		{"никитич", Accusative, "никитича"},
		{"никитич", Instrumental, "никитичем"},
		{"никитич", Prepositional, "никитиче"},

		// -ьич (Ильич) — stem "иль".
		{"ильич", Nominative, "ильич"},
		{"ильич", Genitive, "ильича"},
		{"ильич", Dative, "ильичу"},
		{"ильич", Accusative, "ильича"},
		{"ильич", Instrumental, "ильичем"},
		{"ильич", Prepositional, "ильиче"},
	}
	for _, tc := range cases {
		name := tc.patr + "/" + caseTag(tc.c)
		t.Run(name, func(t *testing.T) {
			got, ok := declinePatronymicHeuristic(tc.patr, tc.c, Masculine)
			if !ok {
				t.Fatalf("declinePatronymicHeuristic(%q, %v, M) — pattern did not match",
					tc.patr, tc.c)
			}
			if got != tc.want {
				t.Fatalf("declinePatronymicHeuristic(%q, %v, M) = %q, expected %q",
					tc.patr, tc.c, got, tc.want)
			}
		})
	}
}

// TestPatronymicHeuristicFemn covers feminine patronymics in -овна, -евна,
// -ична, -инична in six cases.
func TestPatronymicHeuristicFemn(t *testing.T) {
	cases := []patrCase{
		// -овна (Ивановна).
		{"ивановна", Nominative, "ивановна"},
		{"ивановна", Genitive, "ивановны"},
		{"ивановна", Dative, "ивановне"},
		{"ивановна", Accusative, "ивановну"},
		{"ивановна", Instrumental, "ивановной"},
		{"ивановна", Prepositional, "ивановне"},

		// -евна (Сергеевна).
		{"сергеевна", Nominative, "сергеевна"},
		{"сергеевна", Genitive, "сергеевны"},
		{"сергеевна", Dative, "сергеевне"},
		{"сергеевна", Accusative, "сергеевну"},
		{"сергеевна", Instrumental, "сергеевной"},
		{"сергеевна", Prepositional, "сергеевне"},

		// -ична (Никитична) — without a second "и" before "ч".
		{"никитична", Nominative, "никитична"},
		{"никитична", Genitive, "никитичны"},
		{"никитична", Dative, "никитичне"},
		{"никитична", Accusative, "никитичну"},
		{"никитична", Instrumental, "никитичной"},
		{"никитична", Prepositional, "никитичне"},

		// -инична (Ильинична).
		{"ильинична", Nominative, "ильинична"},
		{"ильинична", Genitive, "ильиничны"},
		{"ильинична", Dative, "ильиничне"},
		{"ильинична", Accusative, "ильиничну"},
		{"ильинична", Instrumental, "ильиничной"},
		{"ильинична", Prepositional, "ильиничне"},
	}
	for _, tc := range cases {
		name := tc.patr + "/" + caseTag(tc.c)
		t.Run(name, func(t *testing.T) {
			got, ok := declinePatronymicHeuristic(tc.patr, tc.c, Feminine)
			if !ok {
				t.Fatalf("declinePatronymicHeuristic(%q, %v, F) — pattern did not match",
					tc.patr, tc.c)
			}
			if got != tc.want {
				t.Fatalf("declinePatronymicHeuristic(%q, %v, F) = %q, expected %q",
					tc.patr, tc.c, got, tc.want)
			}
		})
	}
}

// TestPatronymicGenderDetection — suffix-based gender classification.
func TestPatronymicGenderDetection(t *testing.T) {
	cases := []struct {
		patr string
		want Gender
	}{
		{"Иванович", Masculine},
		{"Сергеевич", Masculine},
		{"Никитич", Masculine},
		{"Ильич", Masculine},
		{"Ивановна", Feminine},
		{"Сергеевна", Feminine},
		{"Никитична", Feminine},
		{"Ильинична", Feminine},
		// Uninformative ending — not a patronymic.
		{"Иванов", GenderUnknown},
	}
	for _, tc := range cases {
		t.Run(tc.patr, func(t *testing.T) {
			got := genderFromPatronymic(tc.patr)
			if got != tc.want {
				t.Fatalf("genderFromPatronymic(%q) = %v, expected %v",
					tc.patr, got, tc.want)
			}
		})
	}
}
